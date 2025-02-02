package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"os"
	"path"
)

var (
	LogFile = "$HOME/.dphoto/logs/dphoto.log"
	debug   = false
	Owner   string // Owner source of truce is viper config, for convenience, other commands can get it from here.

	postRunFunctions []func() error

	factory pkgfactory.Factory // factory is set only when the application is ignited (all commands except configure and version)
)

var rootCmd = &cobra.Command{
	Use:   "dphoto",
	Short: "Backup photos and videos to your personal AWS Cloud",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// send all logging to a file to not pollute STDOUT
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

		log.WithFields(log.Fields{
			"LogLevel": log.GetLevel(),
		}).Debugln("Logger setup, starts program...")

		// complete initialisation on components
		if cmd.Name() != "version" {
			ignite := cmd.Name() != "configure"
			factory, err = config.Connect(ignite, cmd.Name() == "configure")
			if err != nil {
				panic(fmt.Errorf("Fatal error while loading configuration: %s \n", err))
			}

			if ignite {
				config.Listen(func(c config.Config) {
					Owner = c.GetString("owner")
				})
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Debugf("Running %d postRunFunction ...", len(postRunFunctions))
		for _, callback := range postRunFunctions {
			err := callback()
			log.Warnf("The %T function failed to complete: %s", callback, err.Error())
		}

		log.Debugln("Program complete.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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
