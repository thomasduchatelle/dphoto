package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	config2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"os"
	"path"
)

var (
	LogFile = "$HOME/.dphoto/logs/dphoto.log"
	debug   = false
	Owner   string // Owner source of truce is viper config, for convenience, other commands can get it from here.

	postRunFunctions []func() error
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
			err = config2.Connect(ignite, cmd.Name() == "configure")
			if err != nil {
				panic(fmt.Errorf("Fatal error while loading configuration: %s \n", err))
			}

			if ignite {
				config2.Listen(func(c config2.Config) {
					Owner = c.GetString("owner")
				})
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Debugf("Running %d postRunFunction ...", len(postRunFunctions))
		for _, callbacks := range postRunFunctions {
			err := callbacks()
			log.Warnf("A function failed to complete: %s", err.Error())
		}

		log.Debugln("Program complete.")
		log.Exit(0)
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
	rootCmd.PersistentFlags().StringVar(&config2.ForcedConfigFile, "config", "", "use configuration file provided instead of searching in ./ , $HOME/.dphoto, and /etc/dphoto")
	rootCmd.PersistentFlags().StringVar(&config2.Environment, "env", "", "add suffix to configuration filename: '--env dev' would use $HOME/dphoto-dev.yml file.")
}
