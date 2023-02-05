package compile

import (
	"fmt"
	"os/exec"
)

// ExecuteCommand executes a command, eventually it will probably do something fancier
func ExecuteCommand(cmd []string) error {
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	fmt.Println(string(out))
	return err

}
