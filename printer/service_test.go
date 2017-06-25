package printer

import (
	"context"
	"image"

	"github.com/go-kit/kit/log"
	. "gopkg.in/check.v1"

	"github.com/fxkr/weblabel/renderer"
)

type ServiceSuite struct {
	ctx      context.Context
	logger   log.Logger
	printer  MockPrinter
	renderer MockRendererService
}

func (s *ServiceSuite) SetUpTest(c *C) {
	s.ctx = context.Background()
	s.logger = log.NewNopLogger()
	s.printer = MockPrinter{}
	s.renderer = MockRendererService{}
}

func (s *ServiceSuite) TestStatus(c *C) {
	service := NewService(&s.printer, &s.renderer, s.logger)
	c.Assert(service.Status(s.ctx), IsNil)
}

func (s *ServiceSuite) TestPrint(c *C) {
	req := printRequest{renderer.Document{Text: "hello"}}
	service := NewService(&s.printer, &s.renderer, s.logger)
	c.Assert(service.Print(s.ctx, req), IsNil)
	c.Assert(len(s.printer.Images), Equals, 1)
}

func (s *ServiceSuite) TestPrintImage(c *C) {
	img := image.NewRGBA(image.Rect(0, 0, 64, 64))
	service := NewService(&s.printer, &s.renderer, s.logger)
	c.Assert(service.PrintImage(s.ctx, img), IsNil)
	c.Assert(len(s.printer.Images), Equals, 1)
	c.Assert(s.printer.Images[0], Equals, img)
}
