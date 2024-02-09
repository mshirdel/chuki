package cmd

import (
	"github.com/mshirdel/chuki/app"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start API server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() {
	app := app.New(configPath)
	err := app.InitAll()
	if err != nil {
		logrus.Fatalf("error in initializing application: %s", err)
	}

	logrus.Info("app started")
}
