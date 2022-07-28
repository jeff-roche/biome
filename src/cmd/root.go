package cmd

import (
	"os"

	"github.com/jeff-roche/biome/src/services"
	"github.com/spf13/cobra"
)

var biomeService *services.BiomeConfigurationService

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "biome",
	Short: "An environment configuration tool",
	Long: `An environment variable configuration tool with capabilities such as:
	- CLI input
	- Configuration commands
	- AWS environment configuration
	- AWS Secrets Manager support
	- Environment variable decryption`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(biomeSvc *services.BiomeConfigurationService) {
	biomeService = biomeSvc

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
