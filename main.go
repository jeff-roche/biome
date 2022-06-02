package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jeff-roche/biome/parser"
)

var version string

func main() {
	vflag := flag.Bool("version", false, "display the release version of this tool")
	biomeName := flag.String("b", "", "[Required] The name of the biome to use")

	flag.Parse()

	// If they specified -v, just print the version
	if *vflag {
		fmt.Println(version)
		return
	}

	// Fetch the command they want to run
	cmds := flag.Args()

	if len(cmds) < 1 {
		log.Fatalln("No command provided")
	}

	// Load the Biome configuration
	err := parser.LoadBiome(*biomeName)
	if err != nil {
		log.Fatalf("unable to load biome '%s': %v", *biomeName, err)
	}

	// Setup and configure the Biome for command execution
	err = parser.ConfigureBiome()
	if err != nil {
		log.Fatalf("unable to configure biome '%s': %v", parser.LoadedBiome.Name, err)
	}

	// Execute order 66
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
