package command

import "os/exec"

func Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	configureHide(cmd)
	return cmd
}
