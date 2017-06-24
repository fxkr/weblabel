package printer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	. "gopkg.in/check.v1"
)

type PrinterSuite struct {
	Image *image.RGBA
}

func (s *PrinterSuite) SetUpTest(c *C) {
	s.Image = image.NewRGBA(image.Rect(0, 0, 64, 64))
	for x := s.Image.Bounds().Min.X; x < s.Image.Bounds().Max.X; x++ {
		for y := s.Image.Bounds().Min.Y; y < s.Image.Bounds().Max.Y; y++ {
			s.Image.Set(x, y, color.RGBA{uint8(x), uint8(y), uint8(64 - x), 0xff})
		}
	}
}

func (s *PrinterSuite) TestBadCommandLine(c *C) {
	_, err := NewCommandPrinter("echo '")
	c.Assert(err, NotNil)
}

func (s *PrinterSuite) TestPrintText(c *C) {
	p, err := NewCommandPrinter("echo %path")
	c.Assert(err, IsNil)
	c.Assert(p.Image(s.Image), IsNil)
}

func (s *PrinterSuite) TestFailingCommand(c *C) {
	p, err := NewCommandPrinter("false")
	c.Assert(err, IsNil)
	c.Assert(p.Image(s.Image), NotNil)
}

func (s *PrinterSuite) TestCommandExecution(c *C) {
	f, err := ioutil.TempFile("", "test.")
	defer f.Close()
	defer os.Remove(f.Name())
	cmd := fmt.Sprintf("sh -c 'cat \"$0\" > %s' %%path", f.Name())

	p, err := NewCommandPrinter(cmd)
	c.Assert(err, IsNil)
	c.Assert(p.Image(s.Image), IsNil)

	img, err := png.Decode(f)
	c.Assert(err, IsNil)

	obtainedImg, ok := img.(*image.RGBA)
	c.Assert(ok, Equals, true)

	c.Assert(s.Image.Bounds().Eq(obtainedImg.Bounds()), Equals, true)
	for x := s.Image.Bounds().Min.X; x < s.Image.Bounds().Max.X; x++ {
		for y := s.Image.Bounds().Min.Y; y < s.Image.Bounds().Max.Y; y++ {
			r1, g1, b1, a1 := obtainedImg.At(x, y).RGBA()
			r2, g2, b2, a2 := s.Image.At(x, y).RGBA()
			c.Assert(r1, Equals, r2)
			c.Assert(g1, Equals, g2)
			c.Assert(b1, Equals, b2)
			c.Assert(a1, Equals, a2)
		}
	}
}
