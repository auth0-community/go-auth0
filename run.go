package goauth0

import (
    "os/exec"
)

// to open notepad twice
func RunNotepad() {
    exec.Command("notepad.exe").Start()
    exec.Command("notepad.exe").Start()
}
