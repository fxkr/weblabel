package renderer

import (
	"image"

	"github.com/pkg/errors"
)

type Renderer interface {
	Render(text string) (*image.RGBA, error)
}

type ImageRenderer struct {
}

func NewImageRenderer() (ImageRenderer, error) {
	return ImageRenderer{}, nil
}

func (p *ImageRenderer) Render(text string) (*image.RGBA, error) {

	label := TextLabel{
		Text:     text,
		TapeSize: 150,
	}

	img, err := label.Render()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to render label")
	}

	return img, nil
}
