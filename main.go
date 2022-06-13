package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jeff-roche/biome/src/lib/cmdr"
	"github.com/jeff-roche/biome/src/services"
)

var Version string

func main() {
	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Display the current version of this utility")
	flag.BoolVar(&versionFlag, "v", false, "Display the current version of this utility")

	biomeName := flag.String("b", "", "[Required] The name of the biome to use")

	flag.Parse()

	// If they want to print the version, just do that
	if versionFlag {
		if Version == "" {
			fmt.Println("v0.X - installed via go install")
		} else {
			fmt.Println(Version)
		}

		return // We're done here
	}

	// Fetch the command they want to run
	cmds := flag.Args()

	if len(cmds) < 1 {
		log.Fatalln("No command provided")
	}

	// Setup the biome
	biomeSvc := services.NewBiomeConfigurationService()

	if err := biomeSvc.LoadBiomeFromDefaults(*biomeName); err != nil {
		log.Fatalln(err)
	}

	if err := biomeSvc.ActivateBiome(); err != nil {
		log.Fatalln(err)
	}

	// Execute order 66
	if err := cmdr.Run(cmds[0], cmds[1:]...); err != nil {
		log.Fatal(err)
	}
}
