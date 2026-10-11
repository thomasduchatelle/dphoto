package cmd

import (
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	LogFile = "$HOME/.dphoto/logs/dphotops.log"
	debug   = false

	factory pkgfactory.Factory
)

var rootCmd = &cobra.Command{
	Use:   "dphotops",
	Short: "Administrative CLI for DPhoto - direct access to AWS resources",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := os.MkdirAll(path.Dir(os.ExpandEnv(LogFile)), 0766)
		if err != nil {
			panic(err)
		}

		openLogFile, err := os.OpenFile(os.ExpandEnv(LogFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err.Error())
		}

		log.SetOutput(openLogFile)
		formatter := new(log.TextFormatter)
		formatter.FullTimestamp = true
		formatter.DisableColors = true
		log.SetFormatter(formatter)
		log.RegisterExitHandler(func() {
			_ = openLogFile.Close()
		})

		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		factory, err = config.Connect(true, false)
		if err != nil {
			panic(fmt.Errorf("Fatal error while loading configuration: %s \n", err))
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
	rootCmd.PersistentFlags().StringVar(&config.ForcedConfigFile, "config", "", "use configuration file provided instead of searching in ./ , $HOME/.dphoto, and /etc/dphoto")
	rootCmd.PersistentFlags().StringVar(&config.Environment, "env", "", "add suffix to configuration filename: '--env dev' would use $HOME/dphoto-dev.yml file.")
}
