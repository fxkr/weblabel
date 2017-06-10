package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fxkr/weblabel/renderer"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

var (
	font  = flag.String("font", "/usr/share/fonts/dejavu/DejaVuSans.ttf", "path to the ttf font")
	text  = flag.String("text", ":-)", "one line of text to put on the label")
	file  = flag.String("file", "label.png", "output filename")
	debug = flag.Bool("debug", false, "enable visual debugging aids")
)

func run(logger log.Logger) error {
	flag.Parse()

	label := renderer.TextLabel{
		Text:     *text,
		Font:     *font,
		TapeSize: 150,
		Debug:    *debug,
	}

	if err := label.RenderFile(*file); err != nil {
		return errors.Wrap(err, "Failed to render label to file")
	}

	return nil
}

func main() {

	// Logging
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	err := run(logger)
	if err != nil {
		logger.Log("component", "main", "err", fmt.Sprintf("%+v", errors.WithStack(err)))
	}
}
