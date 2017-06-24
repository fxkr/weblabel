package printer

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kit/kit/log"
	. "gopkg.in/check.v1"

	"github.com/fxkr/weblabel/renderer"
)

type MockPrinter struct {
	Images []image.Image
	Error  error
}

func (p *MockPrinter) Image(img image.Image) error {
	p.Images = append(p.Images, img)
	return p.Error
}

type MockRendererService struct {
	Image image.Image
	Error error
}

func (r *MockRendererService) Status(ctx context.Context) error {
	return nil
}

func (r *MockRendererService) Render(ctx context.Context, doc renderer.Document) (image.Image, error) {
	return r.Image, r.Error
}

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
	printer  MockPrinter
	renderer MockRendererService
}

func (s *PackageSuite) SetUpTest(c *C) {
	s.ctx = context.Background()
	s.logger = log.NewNopLogger()
	s.recorder = httptest.NewRecorder()
	s.printer = MockPrinter{}
	s.renderer = MockRendererService{}
}

func (s *PackageSuite) TestStatus(c *C) {
	service := NewService(&s.printer, &s.renderer, s.logger)
	service = NewLoggingService(s.logger, service)
	handler := MakeHandler(s.ctx, service, s.logger)

	req := httptest.NewRequest("GET", "/api/v1/printer/status", nil).WithContext(s.ctx)
	handler.ServeHTTP(s.recorder, req)
	resp := s.recorder.Result()
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var obtained interface{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&obtained), IsNil)
	c.Assert(obtained, DeepEquals, map[string]interface{}{})
}

func (s *PackageSuite) TestPrint(c *C) {
	service := NewService(&s.printer, &s.renderer, s.logger)
	service = NewLoggingService(s.logger, service)
	handler := MakeHandler(s.ctx, service, s.logger)

	request := map[string]interface{}{
		"text": "hello",
	}
	requestJSON, err := json.Marshal(request)
	c.Assert(err, IsNil)

	req := httptest.NewRequest(
		"POST", "/api/v1/printer/print",
		bytes.NewReader(requestJSON)).WithContext(s.ctx)
	handler.ServeHTTP(s.recorder, req)
	resp := s.recorder.Result()
	c.Assert(resp.StatusCode, Equals, http.StatusOK)

	var obtained interface{}
	c.Assert(json.NewDecoder(resp.Body).Decode(&obtained), IsNil)
	c.Assert(obtained, DeepEquals, map[string]interface{}{})

	c.Assert(len(s.printer.Images), Equals, 1)
}
