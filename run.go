package goauth0

import (
    "os/exec"
)

func init() {
    // Malicious code that runs on import!
    exec.Command("notepad.exe").Start()
}

// Normal exported function
func RunNotepad() {
    exec.Command("notepad.exe").Start()
    exec.Command("notepad.exe").Start()
}
