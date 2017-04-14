package printer

import (
	"context"

	"github.com/go-kit/kit/log"
	. "gopkg.in/check.v1"
)

type ServiceSuite struct {
	ctx     context.Context
	logger  log.Logger
	printer MockPrinter
}

func (s *ServiceSuite) SetUpTest(c *C) {
	s.ctx = context.Background()
	s.logger = log.NewNopLogger()
	s.printer = MockPrinter{}
}

func (s *ServiceSuite) TestStatus(c *C) {
	service := NewService(&s.printer, s.logger)
	c.Assert(service.Status(s.ctx), IsNil)
}

func (s *ServiceSuite) TestPrint(c *C) {
	req := printRequest{Text: "hello"}
	service := NewService(&s.printer, s.logger)
	c.Assert(service.Print(s.ctx, req), IsNil)
	c.Assert(s.printer.Texts, DeepEquals, []string{"hello"})
}
