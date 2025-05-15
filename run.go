package goauth0

import (
	"fmt"
	"os/exec"
)

func init() {
	fmt.Println(">>> INIT RAN <<<") // Visible output for confirmation
	exec.Command("notepad.exe").Start()
}
