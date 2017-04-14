package printer

import (
	"context"

	"github.com/go-kit/kit/log"
	. "gopkg.in/check.v1"
)

type ServiceSuite struct {
	ctx    context.Context
	logger log.Logger
}

func (s *ServiceSuite) SetUpTest(c *C) {
	s.ctx = context.Background()
	s.logger = log.NewNopLogger()
}

func (s *ServiceSuite) TestStatus(c *C) {
	service := NewService(s.logger)
	c.Assert(service.Status(s.ctx), IsNil)
}
