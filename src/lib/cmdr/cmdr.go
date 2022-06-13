package cmdr

import (
	"os"
	"os/exec"
)

func Run(cmdStr string, args ...string) error {
	cmd := exec.Command(cmdStr, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
