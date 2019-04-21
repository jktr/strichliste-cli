package cmd

import (
	s "github.com/jktr/go-strichliste"
	"github.com/spf13/cobra"
	"os/user"
)

func NewRootCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{

		Use:               "strichliste-cli",
		Short:             "command line interface for strichliste",
		PersistentPreRunE: cli.wrap(initConfig, initClient),
		RunE:              func(cmd *cobra.Command, _ []string) error { return cmd.Usage() },
	}

	cmd.AddCommand(
		newDebitCommand(cli),
		newCreditCommand(cli),
		newRevertCommand(cli),
		newUserCommand(cli),
		newArticleCommand(cli),
		newBuyCommand(cli),
		newMetricsCommand(cli),
		newSettingsCommand(cli),
	)

	cmd.PersistentFlags().String("config", "",
		`config file (default "$XDG_CONFIG_HOME/strichliste-cli/config.json")`)

	user, _ := user.Current()
	cmd.PersistentFlags().StringP("user", "u", user.Username, "your username on strichliste")
	cli.Viper.BindPFlag("user", cmd.PersistentFlags().Lookup("user"))

	cmd.PersistentFlags().String("api-url", "http://[::1]:8080", "strichliste api endpoint")
	cli.Viper.BindPFlag("api-url", cmd.PersistentFlags().Lookup("api-url"))

	return cmd
}

func initConfig(cli *CLI, cmd *cobra.Command, args []string) error {

	configFile, _ := cmd.Flags().GetString("config")
	if configFile != "" {
		cli.Viper.SetConfigFile(configFile)
	} else {
		cli.Viper.SetConfigName("config")
		cli.Viper.AddConfigPath("$XDG_CONFIG_HOME/strichliste-cli/")
		cli.Viper.AddConfigPath("$HOME/.strichliste-cli/")
		cli.Viper.AddConfigPath(".")
	}

	cli.Viper.ReadInConfig()
	return nil
}

func initClient(cli *CLI, cmd *cobra.Command, args []string) error {
	cli.Client = s.NewClient(
		//s.WithApplication("strichliste-cli", "0.1"),
		s.WithEndpoint(cli.Viper.GetString("api-url")),
	)
	return nil
}
