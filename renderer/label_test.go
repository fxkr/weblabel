package renderer

import (
	"image"
	"testing"

	. "gopkg.in/check.v1"
)

var (
	labelTestFont = "/usr/share/fonts/dejavu/DejaVuSans.ttf"
)

func Test(t *testing.T) {
	_ = Suite(&TextLabelSuite{
		Label: &TextLabel{
			Font:     labelTestFont,
			Text:     "N|",
			TapeSize: 100,
		},
		ExpectedEdges: 7, // image ends at right border of "|"
	})
	_ = Suite(&TextLabelSuite{
		Label: &TextLabel{
			Font:     labelTestFont,
			Text:     "_N|_", // ends at border of "_", but it's below center
			TapeSize: 100,
		},
		ExpectedEdges: 8,
	})
	_ = Suite(&TextLabelSuite{
		Label: &TextLabel{
			Font:       labelTestFont,
			Text:       "N|", // image ends after some space after the "|"
			UseAdvance: true,
			TapeSize:   100,
		},
		ExpectedEdges: 8,
	})
	TestingT(t)
}

type TextLabelSuite struct {
	Label         Label
	ExpectedEdges int
}

func (s *TextLabelSuite) TestRender(c *C) {
	img, err := s.Label.Render()
	c.Assert(err, IsNil)
	c.Assert(img.Bounds().Min, DeepEquals, image.Point{0, 0})
	c.Assert(img.Bounds().Dy(), Equals, 100)

	// Check number of edges on horizontal line at center height
	numEdges := 0
	wasBlack := false
	yMiddle := img.Bounds().Max.Y / 2
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		r, g, b, _ := img.At(x, yMiddle).RGBA()
		isBlack := r == 0 && g == 0 && b == 0
		if wasBlack != isBlack {
			numEdges += 1
			wasBlack = isBlack
		}
	}
	c.Assert(numEdges, Equals, s.ExpectedEdges)
}

func (s *TextLabelSuite) TestRenderDebug(c *C) {
	img, err := s.Label.Render()
	c.Assert(err, IsNil)
	c.Assert(img.Bounds().Min, DeepEquals, image.Point{0, 0})
	c.Assert(img.Bounds().Dy(), Equals, 100)
}
