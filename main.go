package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
	dither "github.com/esimov/dithergo"
	v2 "github.com/piger/gattoprint/v2"
	"tinygo.org/x/bluetooth"
)

const (
	PrintWidth = 384
	ErrMul     = 1.18 // error multiplier
)

var ditherers []dither.Dither = []dither.Dither{
	dither.Dither{
		"FloydSteinberg",
		dither.Settings{
			[][]float32{
				[]float32{0.0, 0.0, 0.0, 7.0 / 48.0, 5.0 / 48.0},
				[]float32{3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0},
				[]float32{1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0},
			},
		},
	},
}

func runBluetoothTest() error {
	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}

	err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		fmt.Printf("found device: %s %d %s\n", device.Address.String(), device.RSSI, device.LocalName())
	})
	if err != nil {
		return err
	}

	// GB03
	// found device: 657b44c5-d2b2-69e2-2c52-f33aecfb4a6f -70 GB03

	return nil
}

func blackOrWhite(g color.Gray) color.Gray {
	if g.Y < 127 {
		return color.Gray{0} // Black
	}
	return color.Gray{255} // White
}

func run() error {
	fd, err := os.Open("groucho-marx.jpg")
	if err != nil {
		return err
	}
	defer fd.Close()

	im, imfmt, err := image.Decode(fd)
	if err != nil {
		return err
	}
	bounds := im.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	fmt.Printf("image (%s) size: %dx%d\n", imfmt, width, height)

	factor := float64(PrintWidth) / float64(width)
	newHeight := int(float64(height) * factor)
	fmt.Println("to: ", factor, float64(height)*factor, newHeight)

	dstImage := imaging.Resize(im, PrintWidth, newHeight, imaging.Lanczos)

	// XXX: we should do padding so that the image matches the print width!

	newBounds := dstImage.Bounds()
	gray := image.NewGray(newBounds)
	for x := newBounds.Min.X; x < newBounds.Dx(); x++ {
		for y := newBounds.Min.Y; y < newBounds.Dy(); y++ {
			pixel := dstImage.At(x, y)
			gray.Set(x, y, pixel)
		}
	}

	white := image.NewGray(newBounds)
	for i := range white.Pix {
		white.Pix[i] = 255
	}

	for x := newBounds.Min.X; x < newBounds.Dx(); x++ {
		for y := newBounds.Min.Y; y < newBounds.Dy(); y++ {
			c := blackOrWhite(gray.GrayAt(x, y))
			white.Set(x, y, c)
		}
	}

	ditherer := ditherers[0]
	goo := ditherer.Monochrome(gray, float32(ErrMul))

	b := white.Bounds()
	fmt.Printf("gray image: width=%d, height=%d, stride=%d\n", b.Dx(), b.Dy(), white.Stride)
	fmt.Printf("len pix = %d\n", len(white.Pix))

	// this is how you read an image "line by line"?
	for i := 0; i < len(white.Pix); i += white.Stride {
		row := white.Pix[i : i+white.Stride]
		fmt.Printf("len(row) = %d\n", len(row))
	}

	fmt.Printf("Color Model: %+v\n", white.ColorModel() == color.GrayModel)

	out, err := os.Create("output.png")
	if err != nil {
		return err
	}
	defer out.Close()

	if err := png.Encode(out, goo); err != nil {
		return err
	}

	queue := v2.PrintImage(goo.(*image.Gray))

	if err := v2.SendCommands(queue); err != nil {
		fmt.Printf("error sending commands: %s\n", err)
	}

	// NOTE: the original code "invert" the image using the "~" operator...
	// https://stackoverflow.com/questions/8305199/the-tilde-operator-in-python

	return nil
}

func main() {
	/*
		gh := convert(checksumTable)
		for _, i := range gh {
			fmt.Printf("0x%X ", i)
		}
		fmt.Println()
	*/

	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
