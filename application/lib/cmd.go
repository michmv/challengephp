package lib

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

type CLI struct {
	Use      string
	Short    string
	Long     string
	Commands []Command
}

type Command struct {
	Use   string
	Short string
	Long  string
	Args  cobra.PositionalArgs
	Run   func(cmd *cobra.Command, args []string)
	Init  func(cmd *cobra.Command)
}

func cmdInit(cli CLI) cobra.Command {
	rootCmd := cobra.Command{
		Use:   cli.Use,
		Short: cli.Short,
		Long:  cli.Long,
	}

	for _, command := range cli.Commands {
		subCmd := cobra.Command{
			Use:   command.Use,
			Short: command.Short,
			Long:  command.Long,
			Args:  command.Args,
			Run:   command.Run,
		}
		if command.Init != nil {
			command.Init(&subCmd)
		}
		rootCmd.AddCommand(&subCmd)
	}

	return rootCmd
}

func Start(cli CLI) Error {
	command := cmdInit(cli)
	if err := command.Execute(); err != nil {
		return Err(err)
	}
	return nil
}

func Exit(err Error) {
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
