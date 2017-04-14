package weblabel

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"

	"github.com/fxkr/weblabel/printer"
)

type Config struct {

	// Address and port to listen on, e.g. "0.0.0.0:8000".
	Address string
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
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/weblabel/")
	viper.AddConfigPath("$HOME/.weblabel")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logger.Log("component", "main", "err", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		logger.Log("component", "main", "err", err)
		return
	}

	// Services
	var psLog log.Logger = log.With(logger, "component", "printer")
	var ps printer.Service
	ps = printer.NewService(psLog)
	ps = printer.NewLoggingService(psLog, ps)

	// Service routes
	apiMux := http.NewServeMux()
	apiMux.Handle("/printer/v1/", printer.MakeHandler(ctx, ps, log.With(logger, "component", "http")))

	// Headers
	rootMux := http.NewServeMux()
	rootMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		apiMux.ServeHTTP(w, r)
	}))

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
		logger.Log("component", "main", "err", err)
	case <-sigs:
	}
}
