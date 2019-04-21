package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newBuyCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy",
		Short: "buy some amount of an article",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runBuy),
	}

	cmd.Flags().IntP("article", "a", 0, "id of article to buy")
	cmd.MarkFlagRequired("article")

	cmd.Flags().IntP("count", "c", 1, "amount to buy")

	cmd.Flags().String("comment", "", "add comment to transaction")

	return cmd

}

func runBuy(cli *CLI, cmd *cobra.Command, args []string) error {

	comment, _ := cmd.Flags().GetString("comment")
	articleId, _ := cmd.Flags().GetInt("article")
	count, _ := cmd.Flags().GetInt("count")
	username, _ := cmd.Flags().GetString("user")

	if count <= 0 {
		return fmt.Errorf("must buy at least one instance of the article\n")
	}

	user, _, err := cli.Client.User.GetByName(username)
	if err != nil {
		return err
	}

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	tx, _, err := cli.Client.Transaction.Context(user.ID).
		WithComment(comment).Purchase(articleId, count)
	if err != nil {
		return err
	}

	fmt.Printf("created transaction #%d\n", tx.ID)
	fmt.Printf("new balance for user #%d (%s): %.2f%s\n",
		tx.Issuer.ID,
		tx.Issuer.Name,
		CurrencyIntToFloat64(tx.Issuer.Balance),
		settings.I18n.Currency.Symbol,
	)
	return nil
}
