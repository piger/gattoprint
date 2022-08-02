package graphics

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	dither "github.com/esimov/dithergo"
)

const (
	PrintWidth              = 384  // the width of the printed image, required by the printer.
	errorMultiplier float32 = 1.18 // error multiplier for dithering.
)

var ditherers map[string]dither.Dither

func init() {
	// setup ditherers
	ditherers = make(map[string]dither.Dither)

	fs := dither.Dither{
		Type: "FloydSteinberg",
		Settings: dither.Settings{
			Filter: [][]float32{
				{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
				{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
				{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
			},
		},
	}

	ditherers["FloydSteinberg"] = fs
}

// ConvertImage converts the given image file in a grayscale dithered image,
// scaled down to the width required by the printer.
func ConvertImage(filename string) (*image.Gray, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer fh.Close()

	img, imgFmt, err := image.Decode(fh)
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	log.Printf("decoded image %s as: %s", filename, imgFmt)

	b := img.Bounds()
	width := b.Dx()
	height := b.Dy()
	log.Printf("image size: %dx%d", width, height)

	factor := float64(PrintWidth) / float64(width)
	newHeight := int(float64(height) * factor)

	imgResized := imaging.Fit(img, PrintWidth, newHeight, imaging.Lanczos)
	b = imgResized.Bounds()
	log.Printf("resized to: %dx%d", b.Dx(), b.Dy())

	// XXX: imaging has greyscaling methods as well!
	// imgGray := imaging.Grayscale(imgResized)
	imgGray := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			imgGray.Set(x, y, imgResized.At(x, y))
		}
	}

	ditherer, ok := ditherers["FloydSteinberg"]
	if !ok {
		return nil, fmt.Errorf("ditherer FloydSteinberg not found")
	}

	imgDithered := ditherer.Monochrome(imgGray, errorMultiplier)
	b = imgDithered.Bounds()

	log.Printf("dithered size: %dx%d", b.Dx(), b.Dy())

	if result, ok := imgDithered.(*image.Gray); ok {
		return result, nil
	}

	return nil, errors.New("ditherer did not return the right type")
}
