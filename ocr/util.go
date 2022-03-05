package ocr

import (
	"github.com/disintegration/imaging"
	"image"
	"image/color"
)

type ModifiableImage struct {
	image.Image
}

// At will normalize the c.Y value
func (m *ModifiableImage) At(x, y int) color.Color {
	c := color.Gray16Model.Convert(m.Image.At(x, y)).(color.Gray16)

	if c.Y > 1000 && c.Y < 40000 {
		c.Y += 20000
	}
	return color.Gray16{Y: c.Y}
}

func Run(dir string) error {
	img, err := readImage(dir + "100000.jpg")
	if err != nil {
		return err
	}

	// 13 is the height of the bottom banner. 174 is the width of the logo on the right
	img, err = cropImage(img, image.Rect(0, img.Bounds().Dy()-13, img.Bounds().Dx()-174, img.Bounds().Dy()))
	if err != nil {
		return err
	}

	// Consider using AdjustSigmoid at some point? I didn't really get much further after fiddling with it a ton
	return writeImage(imaging.Invert(&ModifiableImage{img}), dir+"100000-cropped.jpg")
}
