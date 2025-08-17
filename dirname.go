package dirname

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	yup "github.com/yupsh/framework"
	"github.com/yupsh/framework/opt"
	localopt "github.com/yupsh/dirname/opt"
)

// Flags represents the configuration options for the dirname command
type Flags = localopt.Flags

// Command implementation
type command opt.Inputs[string, Flags]

// Dirname creates a new dirname command with the given parameters
func Dirname(parameters ...any) yup.Command {
	return command(opt.Args[string, Flags](parameters...))
}

func (c command) Execute(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer) error {
	if err := yup.RequireArguments(c.Positional, 1, 0, "dirname", stderr); err != nil {
		return err
	}

	separator := "\n"
	if bool(c.Flags.Zero) {
		separator = "\x00"
	}

	for i, path := range c.Positional {
		result := filepath.Dir(path)

		if i > 0 {
			fmt.Fprint(stdout, separator)
		}
		fmt.Fprint(stdout, result)
	}

	if len(c.Positional) > 0 {
		fmt.Fprint(stdout, separator)
	}

	return nil
}

func (c command) String() string {
	return fmt.Sprintf("dirname %v", c.Positional)
}
