package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save the loaded environment variables to a dotenv (.env) file",
	Long: `Save the loaded environment variables to a dotenv (.env) file
	The default file is '.env' in the current directory`,
	Run: func(cmd *cobra.Command, args []string) {
		biomeName, _ := cmd.Flags().GetString("biome")
		fileName, _ := cmd.Flags().GetString("file")
		fmt.Println("LLAMA")
		fmt.Println(biomeName)

		if err := biomeService.LoadBiomeFromDefaults(biomeName); err != nil {
			log.Fatalln(err)
		}

		if err := biomeService.ActivateBiome(); err != nil {
			log.Fatalln(err)
		}

		// Execute order 66
		if err := biomeService.SaveBiomeToFile(fileName); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(saveCmd)

	saveCmd.Flags().StringP("file", "f", ".env", "set the output file name")
	saveCmd.Flags().StringP("biome", "b", "", "the name of the biome to configure")
	saveCmd.MarkFlagRequired("biome")
}
