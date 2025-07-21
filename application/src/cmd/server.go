package cmd

import (
	"challengephp/lib"
	"challengephp/src/di"
	"challengephp/src/response"
	"challengephp/src/routes"
	"challengephp/src/storage"
	"github.com/spf13/cobra"
)

func GetServerCommands(portDefault int, configPathDefault string) lib.Command {
	port := portDefault
	configPath := configPathDefault
	return lib.Command{
		Use:   "server",
		Short: "start the server",
		Run:   runServer(port, configPath),
		Init: func(cmd *cobra.Command) {
			cmd.Flags().IntVarP(&port, "port", "p", port, "start http server on port")
			cmd.Flags().StringVarP(&configPath, "config", "c", configPath, "config file path")
		},
	}
}

func runServer(port int, configPath string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		s, err := storage.New(configPath)
		if err != nil {
			lib.Exit(err.Tap())
		}
		defer s.Close()
		config := lib.ServerConfig[*storage.Storage, *di.Container, response.Response]{
			Storage:            &s,
			RegistrationRoutes: routes.RegistrationRoutes,
			MakeContext:        di.MakeContainer,
			CloseContext:       di.CloseContainer,
		}

		err = lib.RunServer(config, port, s.Debug, s.Logger)
		if err != nil {
			lib.Exit(err.Tap())
		}
	}
}
