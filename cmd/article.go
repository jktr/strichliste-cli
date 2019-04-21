package cmd

import (
	"fmt"
	s "github.com/jktr/go-strichliste"
	"github.com/jktr/go-strichliste/schema"
	"github.com/spf13/cobra"
	"strconv"
)

func newArticleCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "article",
		Short: "interact with the article database",
		Args:  cobra.MaximumNArgs(1),
		RunE:  cli.wrap(runArticleGet),
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create a new article",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runArticleCreate),
	}

	create.Flags().String("name", "", "article's name")
	create.MarkFlagRequired("name")

	create.Flags().Float64("value", 0, "article's value")
	create.MarkFlagRequired("value")

	create.Flags().String("barcode", "", "article's barcode")

	update := &cobra.Command{
		Use:   "update",
		Short: "update an article's metadata",
		Args:  cobra.ExactArgs(1),
		RunE:  cli.wrap(runArticleUpdate),
	}

	update.Flags().String("set-name", "", "article's new name")
	update.Flags().Float64("set-value", 0, "article's new value")
	update.Flags().String("set-barcode", "", "article's new barcode")

	delete := &cobra.Command{
		Use:   "delete",
		Short: "delete/disable an article",
		Args:  cobra.ExactArgs(1),
		RunE:  cli.wrap(runArticleDelete),
	}

	delete.Flags().Bool("confirm", false, "confirm deletion; dry-runs otherwise")

	cmd.AddCommand(create, update, delete)
	return cmd
}

func runArticleGet(cli *CLI, cmd *cobra.Command, args []string) error {

	query := ""
	if len(args) == 1 {
		query = args[0]
	} else {
		return cmd.Usage()
	}

	var articles []schema.Article

	// try interpreting query as ID
	aid, err := strconv.Atoi(query)
	if err == nil {
		article, _, err := cli.Client.Article.Get(aid)
		if err != nil {
			return err
		}
		articles = []schema.Article{*article}
	}

	if len(articles) == 0 {
		articles, _, err = cli.Client.Article.SearchByName(args[0], &s.ListOpts{PerPage: 5})
		if err != nil {
			return err
		}
	}

	if len(articles) == 0 {
		articles, _, err = cli.Client.Article.SearchByBarcode(args[0], &s.ListOpts{PerPage: 5})
		if err != nil {
			return err
		}
	}

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return nil
	}

	for _, article := range articles {
		fmt.Printf("#%03d %s\n", article.ID, article.Name)
		fmt.Printf("\tvalue: %.2f%s\n",
			CurrencyIntToFloat64(article.Value),
			settings.I18n.Currency.Symbol,
		)
		fmt.Printf("\tactive: %t\n", article.IsActive)
		if article.Barcode != nil {
			fmt.Printf("\tbarcode: '%s'\n", *article.Barcode)
		}
	}
	return nil
}

func runArticleCreate(cli *CLI, cmd *cobra.Command, args []string) error {

	name, _ := cmd.Flags().GetString("name")
	floatValue, _ := cmd.Flags().GetFloat64("value")
	barcode, _ := cmd.Flags().GetString("barcode")

	value := CurrencyFloat64ToInt(floatValue)

	article, _, err := cli.Client.Article.Create(&schema.ArticleCreateRequest{
		Name:    name,
		Value:   value,
		Barcode: barcode,
	})
	if err != nil {
		return err
	}

	fmt.Printf("created article #%d (%s)\n", article.ID, article.Name)
	return nil
}

func runArticleUpdate(cli *CLI, cmd *cobra.Command, args []string) error {

	newName, _ := cmd.Flags().GetString("set-name")
	floatValue, _ := cmd.Flags().GetFloat64("set-value")
	newBarcode, _ := cmd.Flags().GetString("set-barcode")

	newValue := CurrencyFloat64ToInt(floatValue)

	articleId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	article, _, err := cli.Client.Article.Get(articleId)
	if err != nil {
		return err
	}

	// don't do anything if nothing changed
	if newName == "" && newBarcode == "" && newValue == 0 {
		fmt.Printf("no updates requested for article #%d (%s)\n", article.ID, article.Name)
		return cmd.Usage()
	}

	// XXX partial update isn't allowed unless the name + value are set,
	// so we set them to their (hopefully still) current value :/

	if newName == "" {
		newName = article.Name
	}
	if newValue == 0 {
		newValue = article.Value
	}

	updatedArticle, _, err := cli.Client.Article.Update(articleId, &schema.ArticleUpdateRequest{
		Name:    newName,
		Value:   newValue,
		Barcode: newBarcode,
	})
	if err != nil {
		return err
	}

	fmt.Printf("update article #%d (%s)\n", updatedArticle.ID, updatedArticle.Name)
	return nil
}

func runArticleDelete(cli *CLI, cmd *cobra.Command, args []string) error {
	articleId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	article, _, err := cli.Client.Article.Get(articleId)
	if err != nil {
		return err
	}

	if !article.IsActive {
		return fmt.Errorf("article is already disabled")
	}

	confirmed, _ := cmd.Flags().GetBool("confirm")
	if !confirmed {
		return fmt.Errorf("dry-run: would delete article #%d (%s)",
			article.ID, article.Name)
	}

	article, _, err = cli.Client.Article.Deactivate(article.ID)
	if err != nil {
		return err
	}

	if article.IsActive {
		return fmt.Errorf("failed to disable article")
	}

	fmt.Printf("disabled article #%d (%s)\n", article.ID, article.Name)
	return nil
}
