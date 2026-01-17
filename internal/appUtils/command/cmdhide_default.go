//go:build !gui

package command

import (
	"os/exec"
)

func configureHide(cmd *exec.Cmd) {}
