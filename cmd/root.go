/*
Copyright © 2020 Brian Williams

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	sourcePath string
	destPath   string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "nginx_flywheel",
		Short: "Configure nginx files with JSON modifications",
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Print(err, "\n")
		os.Exit(1)
	}
}

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	cobra.OnInitialize(initConfig)

	// file flags
	etcdCmd.PersistentFlags().StringVar(&sourcePath, "source", "", "absolute path of NGINX config file")
	etcdCmd.PersistentFlags().StringVar(&destPath, "destination", "", "location to write overwritten values; warning: This will truncate any existing file")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nginx_flywheel.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Print(err, "\n")
			os.Exit(1)
		}

		// Search config in home directory with name ".nginx_flywheel" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".nginx_flywheel")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Print("Using config file:", viper.ConfigFileUsed(), "\n")
	}
}
