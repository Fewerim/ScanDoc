//go:build windows && wails

package command

import (
	"os/exec"
	"syscall"
)

func configureHide(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
