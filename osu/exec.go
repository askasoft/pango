package osu

import (
	"context"
	"errors"
	"io"
	"os/exec"
)

func ExecCommand(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, command string, args ...string) (code int, err error) {
	cmd := BuildCommand(ctx, command, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err = cmd.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			code = ee.ExitCode()
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			err = ctx.Err()
		}
	}

	return
}
