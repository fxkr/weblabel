package printer

import (
	"fmt"
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

type PrinterSuite struct {
}

func (s *PrinterSuite) TestBadCommandLine(c *C) {
	_, err := NewCommandPrinter("echo '")
	c.Assert(err, NotNil)
}

func (s *PrinterSuite) TestPrintText(c *C) {
	p, err := NewCommandPrinter("echo {}")
	c.Assert(err, IsNil)
	c.Assert(p.Text("hello, world"), IsNil)
}

func (s *PrinterSuite) TestFailingCommand(c *C) {
	p, err := NewCommandPrinter("false")
	c.Assert(err, IsNil)
	c.Assert(p.Text("hello, world"), NotNil)
}

func (s *PrinterSuite) TestCommandExecution(c *C) {
	f, err := ioutil.TempFile("", "test.")
	defer f.Close()
	defer os.Remove(f.Name())
	cmd := fmt.Sprintf("sh -c 'echo $0 > %s' {}", f.Name())

	p, err := NewCommandPrinter(cmd)
	c.Assert(err, IsNil)
	c.Assert(p.Text("hello, world"), IsNil)

	obtained, err := ioutil.ReadAll(f)
	c.Assert(err, IsNil)
	c.Assert(string(obtained), Equals, "hello, world\n")
}
