package goauth0

import (
	"os/exec"
)

func init() {
	println("INIT RAN")
	exec.Command("notepad.exe").Start()
}
