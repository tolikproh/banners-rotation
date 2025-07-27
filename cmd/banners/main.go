package main

import (
	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "banners",
	Short: `Cервис "Ротация баннеров"`,
	Long: `Cервис "Ротация баннеров" предназначен для выбора наиболее эффективных (кликабельных) баннеров, 
в условиях меняющихся предпочтений пользователей и набора баннеров.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(version())
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./configs/config.yaml", "path to configuration file")
}

func main() {
	rootCmd.Execute()
}
