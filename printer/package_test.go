package printer

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	_ = Suite(&PackageSuite{})
	_ = Suite(&ServiceSuite{})
	_ = Suite(&PrinterSuite{})
	TestingT(t)
}

type PackageSuite struct {
	ctx      context.Context
	logger   log.Logger
	recorder *httptest.ResponseRecorder
}

func (s *PackageSuite) SetUpTest(c *C) {
	s.ctx = context.Background()
	s.logger = log.NewNopLogger()
	s.recorder = httptest.NewRecorder()
}

func (s *PackageSuite) TestStatus(c *C) {
	service := NewService(s.logger)
	service = NewLoggingService(s.logger, service)
	handler := MakeHandler(s.ctx, service, s.logger)

	req := httptest.NewRequest("GET", "/printer/v1/status", nil).WithContext(s.ctx)
	handler.ServeHTTP(s.recorder, req)
	resp := s.recorder.Result()
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var obtained interface{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&obtained), IsNil)
	c.Assert(obtained, DeepEquals, map[string]interface{}{})
}

func (s *PackageSuite) TestPrint(c *C) {
	service := NewService(s.logger)
	service = NewLoggingService(s.logger, service)
	handler := MakeHandler(s.ctx, service, s.logger)

	request := map[string]interface{}{}
	requestJSON, err := json.Marshal(request)
	c.Assert(err, IsNil)

	req := httptest.NewRequest(
		"POST", "/printer/v1/print",
		bytes.NewReader(requestJSON)).WithContext(s.ctx)
	handler.ServeHTTP(s.recorder, req)
	resp := s.recorder.Result()
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var obtained interface{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&obtained), IsNil)
	c.Assert(obtained, DeepEquals, map[string]interface{}{})
}
