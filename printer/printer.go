package printer

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	"github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
)

type Printer interface {
	Image(img image.Image) error
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

func (p *CommandPrinter) Image(img image.Image) error {
	p.Lock()
	defer p.Unlock()

	f, err := ioutil.TempFile("", "weblabel.")
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.Remove(f.Name())

	if err := png.Encode(f, img); err != nil {
		return errors.Wrap(err, "Failed to encode PNG image")
	}
	f.Close()

	cmd := exec.Command(p.Name, p.getArgs(f.Name())...)
	fmt.Println(p.Name, p.getArgs(f.Name()))
	return cmd.Run()
}

func (p *CommandPrinter) getArgs(path string) []string {
	args := make([]string, len(p.Args))
	copy(args, p.Args)
	for i, _ := range args {
		if args[i] == "%path" {
			args[i] = path
		}
	}
	return args
}
