package cmd

import (
	"log"

	"github.com/jeff-roche/biome/src/lib/cmdr"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run -b <biome-name> [flags] [...cmd]",
	Short: "Run any cli command in the provided biome",
	Long:  "Run any cli command in the provided biome",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		biomeName, _ := cmd.Flags().GetString("biome")

		if err := biomeService.LoadBiomeFromDefaults(biomeName); err != nil {
			log.Fatalln(err)
		}

		if err := biomeService.ActivateBiome(); err != nil {
			log.Fatalln(err)
		}

		// Execute order 66
		if err := cmdr.Run(args[0], args[1:]...); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("biome", "b", "", "the name of the biome to configure")
	runCmd.MarkFlagRequired("biome")
}
