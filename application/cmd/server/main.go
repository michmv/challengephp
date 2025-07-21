package main

import (
	"challengephp/lib"
	"challengephp/src/cmd"
	"fmt"
	"os"
)

func main() {
	configPath := "config.yml"

	cli := lib.CLI{
		Use:   "challengephp",
		Short: "Простой API для сравнения с PHP",
		Commands: []lib.Command{
			cmd.GetServerCommands(12001, configPath),
		},
	}
	err := lib.Start(cli)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
