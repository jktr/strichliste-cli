package cmd

import (
	"fmt"
	s "github.com/jktr/go-strichliste"
	"github.com/jktr/go-strichliste/schema"
	"github.com/spf13/cobra"
	"strconv"
)

func newUserCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "interact with the user database",
		Args:  cobra.MaximumNArgs(1),
		RunE:  cli.wrap(runUserGet),
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "open a new user account",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runUserCreate),
	}

	create.Flags().String("name", "", "user's name")
	create.MarkFlagRequired("name")

	create.Flags().String("email", "", "user's email")
	create.Flags().Float64("balance", 0, "user's initial balance")

	update := &cobra.Command{
		Use:   "update",
		Short: "update a user account's metadata",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runUserUpdate),
	}

	update.Flags().String("set-name", "", "user's new name")
	update.Flags().String("set-email", "", "user's new email")

	delete := &cobra.Command{
		Use:   "delete",
		Short: "delete/disable a user account",
		Args:  cobra.ExactArgs(1),
		RunE:  cli.wrap(runUserDelete),
	}

	delete.Flags().Bool("confirm", false, "confirm deletion; dry-runs otherwise")

	cmd.AddCommand(create, update, delete)
	return cmd
}

func runUserGet(cli *CLI, cmd *cobra.Command, args []string) error {

	query, _ := cmd.Flags().GetString("user")

	// command line has precedence when searching
	if len(args) == 1 {
		query = args[0]
	}

	var users []schema.User

	// try interpreting query as ID
	uid, err := strconv.Atoi(query)
	if err == nil {
		user, _, err := cli.Client.User.Get(uid)
		if err != nil {
			return err
		}
		users = []schema.User{*user}
	}

	// try interpreting query as username
	if len(users) == 0 {
		users, _, err = cli.Client.User.Search(query, &s.ListOpts{PerPage: 5})
		if err != nil {
			return err
		}
	}

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("#%03d %s\n", user.ID, user.Name)
		fmt.Printf("\tbalance: %.2f%s\n",
			CurrencyIntToFloat64(user.Balance),
			settings.I18n.Currency.Symbol,
		)
		fmt.Printf("\tactive: %t\n", user.IsActive)
		if user.Email != nil {
			fmt.Printf("\temail: %s\n", *user.Email)
		}
	}
	return nil
}

func runUserCreate(cli *CLI, cmd *cobra.Command, args []string) error {

	username, _ := cmd.Flags().GetString("name")
	email, _ := cmd.Flags().GetString("email")
	floatBalance, _ := cmd.Flags().GetFloat64("balance")

	balance := CurrencyFloat64ToInt(floatBalance)

	user, _, err := cli.Client.User.Create(&schema.UserCreateRequest{
		Name:  username,
		Email: email,
	})
	if err != nil {
		return err
	}

	fmt.Printf("created user #%d (%s)\n", user.ID, user.Name)

	// fake an inital balance by issuing a transaction
	if balance != 0 {
		settings, _, err := cli.Client.Settings.Get()
		if err != nil {
			return err
		}

		tx, _, err := cli.Client.Transaction.Context(user.ID).
			WithComment("initial balance").Delta(balance)
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
	}

	return nil
}

func runUserUpdate(cli *CLI, cmd *cobra.Command, args []string) error {

	newUsername, _ := cmd.Flags().GetString("set-name")
	newEmail, _ := cmd.Flags().GetString("set-email")
	username, _ := cmd.Flags().GetString("user")

	user, _, err := cli.Client.User.GetByName(username)
	if err != nil {
		return err
	}

	// don't do anything if nothing changed
	if newUsername == "" && newEmail == "" {
		fmt.Printf("no updates requested for user #%d (%s)\n", user.ID, user.Name)
		return cmd.Usage()
	}

	user, _, err = cli.Client.User.Update(user.ID, &schema.UserUpdateRequest{
		Name:  newUsername,
		Email: newEmail,
	})
	if err != nil {
		return err
	}

	fmt.Printf("updated user #%d (%s)\n", user.ID, user.Name)
	return nil
}

func runUserDelete(cli *CLI, cmd *cobra.Command, args []string) error {

	uid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	user, _, err := cli.Client.User.Get(uid)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return fmt.Errorf("user is already disabled")
	}

	confirmed, _ := cmd.Flags().GetBool("confirm")
	if !confirmed {
		return fmt.Errorf("dry-run: would delete user #%d (%s)",
			user.ID, user.Name)
	}

	user, _, err = cli.Client.User.Deactivate(user.ID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return fmt.Errorf("failed to disable user")
	}

	fmt.Printf("disabled user #%d (%s)\n", user.ID, user.Name)
	return nil
}
