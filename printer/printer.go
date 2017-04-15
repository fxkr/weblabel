package printer

import (
	"os/exec"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

type Printer interface {
	Text(text string) error
}

type CommandPrinter struct {
	Name string
	Args []string
}

func NewCommandPrinter(commandLine string) (CommandPrinter, error) {
	args, err := shellwords.Parse(commandLine)
	if err != nil {
		return CommandPrinter{}, errors.WithStack(err)
	}
	return CommandPrinter{args[0], args[1:]}, nil
}

func (p *CommandPrinter) Text(text string) error {
	cmd := exec.Command(p.Name, p.getArgs(text)...)
	return cmd.Run()
}

func (p *CommandPrinter) getArgs(text string) []string {
	args := make([]string, len(p.Args))
	copy(args, p.Args)
	for i, _ := range args {
		if args[i] == "{}" {
			args[i] = text
		}
	}
	return args
}
