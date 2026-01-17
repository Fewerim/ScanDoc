//go:build !(windows && wails)

package command

import (
	"os/exec"
)

func configureHide(cmd *exec.Cmd) {}
