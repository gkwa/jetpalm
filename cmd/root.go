/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	defaultConfigFilename      = "jetpalm"
	envPrefix                  = "STING"
	replaceHyphenWithCamelCase = false
	cfgFile                    string
)

func initializeConfig(cmd *cobra.Command, cfg *Config) error {
	v := createViperInstance()

	SetDefaultConfigValues(cfg)
	// Check if the default config file exists
	if err := WriteDefaultConfigToFile(cfg, defaultConfigFilename+".yaml"); err != nil {
		return err
	}

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.AutomaticEnv()
	bindFlags(cmd, v)

	if err := v.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func getConfigName(f *pflag.Flag) string {
	configName := f.Name
	if replaceHyphenWithCamelCase {
		configName = strings.ReplaceAll(f.Name, "-", "")
	}
	return configName
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := getConfigName(f)
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func createViperInstance() *viper.Viper {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return v
}

func NewRootCommand() *cobra.Command {
	var cfg Config

	rootCmd := &cobra.Command{
		Use:   "jetpalm",
		Short: "Cober and Viper together at last",
		Long:  `Demonstrate how to get cobra flags to bind to viper properly`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd, &cfg)
		},
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			PrintConfigValues(out, &cfg)
		},
	}

	return rootCmd
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCommand()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jetpalm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".jetpalm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jetpalm")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func PrintConfigValues(out io.Writer, cfg *Config) {
	fmt.Fprintln(out, "pushfrequency:", cfg.Client.PushFrequency)
}
