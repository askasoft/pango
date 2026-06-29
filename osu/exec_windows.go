package osu

import (
	"context"
	"os/exec"

	"github.com/askasoft/pango/num"
)

func BuildCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)

	// Tell Go to kill the entire process group immediately on timeout.
	// This prevents Go from hanging on open stdout/stderr pipes.
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			// /F = Force, /T = terminate child processes (tree)
			return exec.Command("taskkill", "/F", "/T", "/PID", num.Itoa(cmd.Process.Pid)).Run()
		}
		return nil
	}

	return cmd
}
