package command

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	yup "github.com/gloo-foo/framework"
)

type command yup.Inputs[string, flags]

func Dirname(parameters ...any) yup.Command {
	return command(yup.Initialize[string, flags](parameters...))
}

func (p command) Executor() yup.CommandExecutor {
	return func(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
		// Process each positional argument
		for _, path := range p.Positional {
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
