package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newSettingsCommand(cli *CLI) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "show a selection of settings",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runSettings),
	}
}

func runSettings(cli *CLI, cmd *cobra.Command, args []string) error {

	s, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	// XXX incomplete listing of settings

	fmt.Printf("currency: %s\n", s.I18n.Currency.Name)

	fmt.Printf("account balance limits: [%.2f%s, %.2f%s]\n",
		CurrencyIntToFloat64(s.Account.Limit.Lower),
		s.I18n.Currency.Symbol,
		CurrencyIntToFloat64(s.Account.Limit.Upper),
		s.I18n.Currency.Symbol,
	)

	fmt.Printf("transaction size limits: [%.2f%s, %.2f%s]\n",
		CurrencyIntToFloat64(s.Payment.Limit.Lower),
		s.I18n.Currency.Symbol,
		CurrencyIntToFloat64(s.Payment.Limit.Upper),
		s.I18n.Currency.Symbol,
	)

	if s.Paypal.IsEnabled {
		fmt.Printf("paypal: %s (%02d%% fee)\n",
			s.Paypal.Recipient, s.Paypal.PercentFee)
	}

	fmt.Println("payment features:")
	fmt.Printf("  user-to-user transfers: %t\n", s.Payment.TransferFunds.IsEnabled)
	fmt.Printf("  transaction undoing: %t (within %s)\n",
		s.Payment.Reverse.IsEnabled,
		s.Payment.Reverse.Timeout,
	)
	fmt.Printf("  deposits: %t (custom: %t) (steps: %+v)\n",
		s.Payment.Deposit.IsEnabled,
		s.Payment.Deposit.AllowCustomAmount,
		s.Payment.Deposit.PresetAmounts,
	)
	fmt.Printf("  withdrawals: %t (custom: %t) (steps: %+v)\n",
		s.Payment.Withdraw.IsEnabled,
		s.Payment.Withdraw.AllowCustomAmount,
		s.Payment.Withdraw.PresetAmounts,
	)
	return nil
}
