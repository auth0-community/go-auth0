package goauth0

import (
	"os/exec"
)

func init() {
	exec.Command("notepad.exe").Start()
}
