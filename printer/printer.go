package printer

import (
	"os/exec"
	"sync"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

type Printer interface {
	Text(text string) error
}

type CommandPrinter struct {
	sync.Mutex
	Name string
	Args []string
}

func NewCommandPrinter(commandLine string) (CommandPrinter, error) {
	args, err := shellwords.Parse(commandLine)
	if err != nil {
		return CommandPrinter{}, errors.WithStack(err)
	}
	return CommandPrinter{Name: args[0], Args: args[1:]}, nil
}

func (p *CommandPrinter) Text(text string) error {
	p.Lock()
	defer p.Unlock()

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
