package ocr

import "C"
import (
	"bytes"
	"github.com/disintegration/imaging"
	"github.com/otiai10/gosseract/v2"
	"image"
	"image/color"
	"image/png"
	"log"
	"ranet-clone/cfg"
	"ranet-clone/threads"
)

var (
	client gosseract.Client
)

func InitClient(tessDataPrefix string, nThreads int) error {
	client = *gosseract.NewClient()
	client.TessdataPrefix = tessDataPrefix
	err := client.SetLanguage([]string{"eng"}...)

	if err != nil {
		return err
	}

	return nil
}

func GeneratePaths(dir string) ([]cfg.ImageInfo, error) {
	files := threads.GetFiles(dir)
	arr := make([]cfg.ImageInfo, 0)

	for _, f := range files {
		if !f.IsDir() {
			arr = append(arr, cfg.ImageInfo{Name: f.Name()})
		}
	}

	return arr, nil
}

func ProcessImages(p []cfg.ImageInfo, dir string) {
	defer threads.LogPanic()

	log.Printf("will process %v images\n", len(p))
	for _, i := range p {
		if cfg.InOcrQueue(i) || len(i.OcrData) > 0 {
			continue
		}

		cfg.AddToOcrQueue(i)
		if str, err := ProcessImage(dir, i.Name); err != nil {
			log.Printf("error processing %s: %v\n", i.Name, err)
		} else {
			cfg.UpdateOcrData(i, str)
		}

		cfg.RemoveFromOcrQueue(i)
	}
}

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

func ProcessImage(dir, name string) (string, error) {
	data, err := GetImageBytes(dir, name)
	if err != nil {
		return "", err
	}

	err = client.SetImageFromBytes(data)
	if err != nil {
		return "", err
	}

	return client.Text()
}

func GetImageBytes(dir, name string) ([]byte, error) {
	img, err := readImage(dir + name)
	if err != nil {
		return nil, err
	}

	// 13 is the height of the bottom banner. 174 is the width of the logo on the right
	img, err = cropImage(img, image.Rect(0, img.Bounds().Dy()-13, img.Bounds().Dx()-174, img.Bounds().Dy()))
	if err != nil {
		return nil, err
	}

	// Consider using AdjustSigmoid at some point? I didn't really get much further after fiddling with it a ton
	img = imaging.Invert(&ModifiableImage{img})

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
