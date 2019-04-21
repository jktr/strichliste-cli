package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func newRevertCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "revert",
		Aliases: []string{"undo"},
		Short:   "delete/reverse a transaction",
		Args:    cobra.ExactArgs(1),
		RunE:    cli.wrap(runRevert),
	}

	cmd.Flags().Bool("confirm", false, "confirm deletion; dry-runs otherwise")

	return cmd
}

func runRevert(cli *CLI, cmd *cobra.Command, args []string) error {

	username, _ := cmd.Flags().GetString("user")
	txId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	user, _, err := cli.Client.User.GetByName(username)
	if err != nil {
		return err
	}

	// XXX any valid user works, for some reason
	context := cli.Client.Transaction.Context(user.ID)

	tx, _, err := context.Get(txId)
	if err != nil {
		return err
	}

	if tx.IsReversed {
		return fmt.Errorf("transaction is already reversed")
	}

	if !tx.IsReversible {
		return fmt.Errorf("transaction cannot be reversed")
	}

	confirmed, _ := cmd.Flags().GetBool("confirm")
	if !confirmed {
		return fmt.Errorf("dry-run: would delete tx #%d with user #%d (%s)",
			tx.ID, user.ID, user.Name)
	}

	rev, _, err := context.Revert(txId)
	if !rev.IsReversed {
		return fmt.Errorf("failed to reverse transaction")
	}

	fmt.Printf("reversed transaction #%d\n", rev.ID)
	return nil

}
