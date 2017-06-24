package weblabel

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/fxkr/weblabel/printer"
	"github.com/fxkr/weblabel/renderer"
)

type Config struct {

	// Address and port to listen on, e.g. "0.0.0.0:8000".
	Address string

	// Command to execute to print. "{}" will be replaced with text.
	PrintCommand string

	// Directory where the static files are. "" to disable serving them.
	StaticPath string
}

func Run() {
	ctx := context.Background()

	// Logging
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Config
	var config Config
	viper.SetDefault("Address", "127.0.0.1:8080")
	viper.SetDefault("PrintCommand", "echo {}")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.weblabel")
	viper.AddConfigPath("/etc/weblabel/")
	viper.AddConfigPath("/usr/share/weblabel/")
	if err := viper.ReadInConfig(); err != nil {
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
		return
	}

	// Printer
	p, err := printer.NewCommandPrinter(config.PrintCommand)
	if err != nil {
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
		return
	}

	// Renderer
	r, err := renderer.NewImageRenderer()
	if err != nil {
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
		return
	}

	// Services
	var rsLog log.Logger = log.With(logger, "component", "renderer")
	var rs renderer.Service
	rs = renderer.NewService(&r, rsLog)
	rs = renderer.NewLoggingService(rsLog, rs)
	var psLog log.Logger = log.With(logger, "component", "printer")
	var ps printer.Service
	ps = printer.NewService(&p, rs, psLog)
	ps = printer.NewLoggingService(psLog, ps)

	// Service routes
	apiMux := http.NewServeMux()
	apiMux.Handle("/api/v1/printer/", printer.MakeHandler(ctx, ps, log.With(logger, "component", "http")))
	apiMux.Handle("/api/v1/renderer/", renderer.MakeHandler(ctx, rs, log.With(logger, "component", "http")))

	// Distinguish API / static files, set appropriate headers
	rootMux := http.NewServeMux()
	rootHandler := makeApiRootHandler(apiMux)
	if config.StaticPath != "" {
		fileServer := http.FileServer(http.Dir(config.StaticPath))
		rootHandler = makeCombinedRootHandler(rootHandler, fileServer)
	}
	rootMux.Handle("/", rootHandler)

	// Serve
	errs := make(chan error, 1)
	go func() {
		errs <- http.ListenAndServe(config.Address, rootMux)
	}()

	// Await termination
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)
	select {
	case err := <-errs:
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
	case <-sigs:
	}
}

func makeCombinedRootHandler(apiHandler, fileHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			apiHandler.ServeHTTP(w, r)
		} else {
			fileHandler.ServeHTTP(w, r)
		}
	}
}

func makeApiRootHandler(apiHandler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		apiHandler.ServeHTTP(w, r)
	}
}
