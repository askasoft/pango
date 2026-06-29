package osu

import (
	"context"
	"os/exec"
	"syscall"
)

func BuildCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)

	// Set Process Group ID so all child processes share the same group
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Tell Go to kill the entire process group immediately on timeout.
	// This prevents Go from hanging on open stdout/stderr pipes.
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			// Negative PID targets the entire process group
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil
	}

	return cmd
}
