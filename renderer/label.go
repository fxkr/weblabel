package renderer

import (
	_ "github.com/golang/freetype"
)

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/fxkr/go-freetype-fontloader"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	DEFAULT_FONT = "DejaVuSans"
)

type Label interface {
	Render() (*image.RGBA, error)
}

type TextLabel struct {
	Text string
	Font string

	// Enable visual debug output.
	Debug bool

	// Width of tape in pixel, corresponds to text height.
	TapeSize int

	// Include space after text in width.
	UseAdvance bool
}

func (l *TextLabel) getFontDrawer(fontSize float64) (*font.Drawer, error) {

	fontName := l.Font
	if fontName == "" {
		fontName = DEFAULT_FONT
	}

	fontFace, err := fontloader.LoadCache(fontName)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load font")
	}

	fontDrawer := &font.Drawer{
		Dst: nil,
		Src: image.Black,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}

	return fontDrawer, nil
}

func (l *TextLabel) Render() (*image.RGBA, error) {

	// Reserve "some" space for ascent, descent and as margins
	fontSize := 0.8 * float64(l.TapeSize)

	// Load font
	fontDrawer, err := l.getFontDrawer(fontSize)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load font")
	}

	// Width of label
	bounds, labelWidth := font.BoundString(fontDrawer.Face, l.Text)
	if !l.UseAdvance {
		labelWidth = bounds.Max.X
	}

	// Create image of appropriate size
	imageHeight := l.TapeSize
	imageWidth := (labelWidth).Round()
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	// Calculate texts origin (its on the baseline)
	metrics := fontDrawer.Face.Metrics()
	dotX := fixed.I(0)
	dotY := fixed.I(l.TapeSize)/2 + (metrics.Ascent-metrics.Descent)/2

	// Draw string
	fontDrawer.Dst = img
	fontDrawer.Dot = fixed.Point26_6{
		X: dotX,
		Y: dotY,
	}
	fontDrawer.DrawString(l.Text)

	// Guide lines for debugging
	if l.Debug {
		for x := 0; x < imageWidth; x++ {
			fg := color.RGBA{255, 0, 255, 255}
			img.Set(x, (dotY - metrics.Ascent).Round(), fg)                   // Ascent
			img.Set(x, l.TapeSize/2, fg)                                      // Tape center
			img.Set(x, dotY.Round(), fg)                                      // Baseline
			img.Set(x, (dotY + metrics.Descent).Round(), fg)                  // Descent
			img.Set(x, (dotY + metrics.Descent - metrics.Height).Round(), fg) // Cap
		}
	}

	return img, nil

}

func (l *TextLabel) RenderFile(path string) error {
	rgba, err := l.Render()
	if err != nil {
		return errors.Wrap(err, "Failed to render label")
	}

	outFile, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "Failed to create output file")
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		return errors.Wrap(err, "Failed to encode PNG image")
	}

	err = b.Flush()
	if err != nil {
		return errors.Wrap(err, "Failed to flush output file")
	}

	return nil
}
