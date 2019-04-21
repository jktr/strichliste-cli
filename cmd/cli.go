package cmd

import (
	s "github.com/jktr/go-strichliste"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math"
)

type CLI struct {
	RootCommand *cobra.Command
	Viper       *viper.Viper
	Client      *s.Client
}

func NewCLI() *CLI {
	cli := &CLI{
		Viper: viper.New(),
	}
	cli.RootCommand = NewRootCommand(cli)
	return cli
}

func (c *CLI) wrap(fs ...func(*CLI, *cobra.Command, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range fs {
			err := f(c, cmd, args)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func CurrencyIntToFloat64(balance int) float64 {
	return float64(balance) / 100
}

func CurrencyFloat64ToInt(balance float64) int {
	return int(math.Round(balance * 100))
}
