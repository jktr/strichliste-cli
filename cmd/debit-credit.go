package cmd

import (
	"fmt"
	"github.com/jktr/go-strichliste/schema"
	"github.com/spf13/cobra"
)

func newDebitCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "debit",
		Aliases: []string{"withdraw"},
		Short:   "deduct from an account's balance",
		Args:    cobra.NoArgs,
		RunE:    cli.wrap(runTransact),
	}

	cmd.Flags().String("from", "", "account to debit (prefer --user)")
	cmd.Flags().String("to", "", "account to credit (if any)")

	cmd.Flags().StringP("comment", "c", "", "add comment to transaction")

	cmd.Flags().Float64P("amount", "a", 0, "amount to withdraw, as a decimal")
	cmd.MarkFlagRequired("amount")

	return cmd

}
func newCreditCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "credit",
		Aliases: []string{"deposit"},
		Short:   "add to an account's balance",
		Args:    cobra.NoArgs,
		RunE:    cli.wrap(runTransact),
	}

	cmd.Flags().String("from", "", "account to debit (if any)")
	cmd.Flags().String("to", "", "account to credit (prefer --user)")

	cmd.Flags().StringP("comment", "c", "", "add comment to transaction")

	cmd.Flags().Float64P("amount", "a", 1.00, "amount to deposit, as a decimal")
	cmd.MarkFlagRequired("amount")

	return cmd

}

func resolveSrcDst(cli *CLI, cmd *cobra.Command) (string, string) {

	// XXX this seems unnecessarily convoluted, but seems to
	// produce the most intuitive command line usage
	//
	// cd
	//   > $user
	// cd from to
	//   $from > $to
	// cd from
	//   $from > $user
	// cd to
	//   > $to
	// dw
	//   $user >
	// dw from to
	//   $from > $to
	// dw from
	//   $from >
	// dw to
	//   $user > $to

	local, _ := cmd.Flags().GetString("user")
	from, _ := cmd.Flags().GetString("from")
	to, _ := cmd.Flags().GetString("to")

	src, dst := "", ""

	switch cmd.Use {
	case "deposit":
		fallthrough
	case "credit":
		dst = local

	case "withdraw":
		fallthrough
	case "debit":
		src = local
	}

	if from != "" {
		src = from
	}

	if to != "" {
		dst = to
	}
	return src, dst
}

func runTransact(cli *CLI, cmd *cobra.Command, args []string) error {

	comment, _ := cmd.Flags().GetString("comment")

	floatAmount, _ := cmd.Flags().GetFloat64("amount")
	amount := CurrencyFloat64ToInt(floatAmount)
	if amount == 0 {
		return fmt.Errorf("amount most not be zero\n")
	}

	src, dst := resolveSrcDst(cli, cmd)

	if src == dst {
		return fmt.Errorf("source and destination must be different when sending funds\n")
	}

	// XXX User-to-User transaction are forced to use negative amounts

	if (src != "" && dst != "") || cmd.Use == "withdraw" || cmd.Use == "debit" {
		amount = -amount
	}

	var srcUser, dstUser *schema.User
	var err error

	if src != "" {
		srcUser, _, err = cli.Client.User.GetByName(src)
		if err != nil {
			return err
		}

	}
	if dst != "" {
		dstUser, _, err = cli.Client.User.GetByName(dst)
		if err != nil {
			return err
		}
	}

	if src == "" {
		return transactDelta(cli, dstUser, amount, comment)
	}

	if dst == "" {
		return transactDelta(cli, srcUser, amount, comment)
	}

	return transactSend(cli, srcUser, dstUser, amount, comment)
}

func transactDelta(cli *CLI, from *schema.User, amount int, comment string) error {

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	tx, _, err := cli.Client.Transaction.Context(from.ID).
		WithComment(comment).Delta(amount)
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

func transactSend(cli *CLI, from, to *schema.User, amount int, comment string) error {

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	tx, _, err := cli.Client.Transaction.Context(from.ID).
		WithComment(comment).TransferFunds(to.ID, amount)
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
	fmt.Printf("new balance for user #%d (%s): %.2f%s\n",
		tx.To.ID,
		tx.To.Name,
		CurrencyIntToFloat64(tx.To.Balance),
		settings.I18n.Currency.Symbol,
	)

	return nil
}
