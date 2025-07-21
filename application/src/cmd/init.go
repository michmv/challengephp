package cmd

import (
	"challengephp/lib"
	"challengephp/src/config"
	"challengephp/src/db"
	"context"
	"fmt"
	"github.com/spf13/cobra"
)

func GetDatabaseInitCommand(configPathDefault string) lib.Command {
	configPath := configPathDefault
	return lib.Command{
		Use:   "init",
		Short: "initialize the database",
		Run:   runInit(configPath),
	}
}

func runInit(configPath string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		conf, err := config.LoadConfig(configPath)
		if err != nil {
			lib.Exit(err.Tap())
		}

		pool, err := db.CreateDB(conf.DB)
		if err != nil {
			lib.Exit(err.Tap())
		}
		defer pool.Close()

		for _, q := range sliceStringMerge(db.Init, db.CreateIndex) {
			_, err2 := pool.Exec(context.Background(), q)
			if err2 != nil {
				lib.Exit(lib.Err(err2))
			}
		}
		fmt.Println("База данных инициализированна")
	}
}
