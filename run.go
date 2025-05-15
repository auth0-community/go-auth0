package goauth0

import (
	"fmt"
	"os/exec"
)

func init() {
	fmt.Println(">>> INIT FUNCTION EXECUTED <<<")
	exec.Command("notepad.exe").Start()
}
