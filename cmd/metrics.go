package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"sort"
)

func newMetricsCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "show a selection of metrics",
		Args:  cobra.NoArgs,
		RunE:  cli.wrap(runMetrics),
	}

	cmd.Flags().Bool("system", false, "show system metrics")

	return cmd
}

func runMetrics(cli *CLI, cmd *cobra.Command, args []string) error {

	system, _ := cmd.Flags().GetBool("system")
	username, _ := cmd.Flags().GetString("user")

	if system {
		return systemMetrics(cli)
	} else {
		return userMetrics(cli, username)
	}
}

func userMetrics(cli *CLI, username string) error {

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	user, _, err := cli.Client.User.GetByName(username)
	if err != nil {
		return err
	}

	m, _, err := cli.Client.Metrics.ForUser(user.ID)
	if err != nil {
		return err
	}

	// XXX incomplete listing of metrics
	// TODO statsd? influx?

	fmt.Printf("current user balance: %.2f%s\n",
		CurrencyIntToFloat64(m.Balance),
		settings.I18n.Currency.Symbol,
	)
	fmt.Printf("total number of transactions: %d\n", m.Transactions.Count)
	fmt.Printf("total funds sent to other users: %d%s\n",
		m.Transactions.Outgoing.Cashflow,
		settings.I18n.Currency.Symbol,
	)
	fmt.Printf("total funds received from other users: %d%s\n",
		m.Transactions.Incoming.Cashflow,
		settings.I18n.Currency.Symbol,
	)

	if len(m.Articles) > 0 {
		sort.Slice(m.Articles, func(i, j int) bool {
			return m.Articles[i].Count < m.Articles[j].Count
		})

		fmt.Println("user's most popular articles:")
		for _, a := range m.Articles {
			fmt.Printf("\t%3d x %s ~= %.2f%s\n",
				a.Count,
				a.Article.Name,
				CurrencyIntToFloat64(a.Spent),
				settings.I18n.Currency.Symbol,
			)
		}

	}

	return nil
}

func systemMetrics(cli *CLI) error {

	settings, _, err := cli.Client.Settings.Get()
	if err != nil {
		return err
	}

	m, _, err := cli.Client.Metrics.ForSystem()
	if err != nil {
		return err
	}

	// XXX incomplete listing of metrics
	// TODO statsd? influx?

	fmt.Printf("current system balance: %.2f%s\n",
		CurrencyIntToFloat64(m.Balance),
		settings.I18n.Currency.Symbol,
	)
	fmt.Printf("total number of transactions: %d\n", m.Transactions)
	fmt.Printf("total number of users: %d\n", m.Users)

	return nil
}
