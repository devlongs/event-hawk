package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CfgFile      string
	OutputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "eth-monitor",
	Short: "Real-time Ethereum event monitor",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&CfgFile, "config", "config.yaml", "Config file path")
	rootCmd.PersistentFlags().StringVarP(&OutputFormat, "output", "o", "text", "Output format (text/json)")
}

func initConfig() {
	viper.SetConfigFile(CfgFile)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error reading config file: %w", err))
	}
}
