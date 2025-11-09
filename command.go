package command

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	gloo "github.com/gloo-foo/framework"
)

type command gloo.Inputs[string, flags]

func Dirname(parameters ...any) gloo.Command {
	return command(gloo.Initialize[string, flags](parameters...))
}

func (p command) Executor() gloo.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Process each positional argument
		for _, path := range p.Positional {
			// Clean path first to handle trailing slashes properly
			path = filepath.Clean(path)
			// Get directory part
			dir := filepath.Dir(path)

			// Use zero as line separator if flag is set
			if bool(p.Flags.Zero) {
				_, err := fmt.Fprintf(stdout, "%s\x00", dir)
				if err != nil {
					return err
				}
			} else {
				_, err := fmt.Fprintln(stdout, dir)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
