package commands

import (
	"image"

	"golang.org/x/exp/constraints"
)

// see sources/com/blueUtils/BluetoothOrder.java

var (
	cmdSetQuality     byte = 0xA4
	cmdControlLattice byte = 0xA6
	cmdSetEnergy      byte = 0xAF
	cmdDrawingMode    byte = 0xBE // 1 for text, 0 for images
	cmdOtherFeedPaper byte = 0xBD
	cmdDrawBitmap     byte = 0xA2 // Line to draw. 0 bit -> don't draw pixel, 1 bit -> draw pixel
	cmdFeedPaper      byte = 0xA1
	cmdGetDevState    byte = 0xA3

	cmdPrintLattice  []byte = []byte{0xAA, 0x55, 0x17, 0x38, 0x44, 0x5F, 0x5F, 0x5F, 0x44, 0x38, 0x2C}
	cmdImgPrintSpeed []byte = []byte{0x23}
	cmdFinishLattice []byte = []byte{0xAA, 0x55, 0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17}
	cmdBlankSpeed    []byte = []byte{0x19}

	// QUALITIES
	// 0x31 0x32 0x33 0x34 0x35 (49, 50, 51, 52, 53 in Java)

	// PRINT_COMMANDS
	// 0xA2 0xBF (fixed length, run-length)

	// CMD_SET_PAPER (which is cmdFeedPaper)
	// [0x51 0x78] 0xA1 0x0 0x2 0x0 0x30 0x0 0xF9 0xFF

	// CMD_SET_QUALITY_200_DPI
	// 0x51 0x78 0xA4 0x0 0x1 0x0 0x32 0x9E 0xFF
	// ^^^^ ^^^^ prefix
	//           ^^^^ command
	//                ^^^ empty
	//                    ^^^ length
	//                            ^^^^ data
	//                                 ^^^^ checksum
	//                                      ^^^^ closing
	// cmd200DPI []byte = []byte{0x01, 0x00, 0x32, 0x9E, 0xFF}
)

func formatMessage(command byte, data []byte) []byte {
	var result []byte = []byte{
		0x51, 0x78, command, 0x00, byte(len(data)), 0x00,
	}
	result = append(result, data...)
	crc := crc8(data)
	result = append(result, crc)
	result = append(result, 0xFF)
	return result
}

func printerShort(i int) []byte {
	var result []byte = []byte{
		byte(i & 0xFF), byte((i >> 8) & 0xFF),
	}
	return result
}

/* contrast

energy = {
    0: printer_short(8000),
    1: printer_short(12000),
    2: printer_short(17500)
}
contrast = 1
*/

func PrintImage(img *image.Gray) [][]byte {
	var queue [][]byte

	// set quality to standard
	c1 := formatMessage(cmdSetQuality, []byte{0x33})
	queue = append(queue, c1)

	// start and/or set up the lattice, whatever that is
	c2 := formatMessage(cmdControlLattice, cmdPrintLattice)
	queue = append(queue, c2)

	// Set energy used
	var contrast int = 12000
	c3 := formatMessage(cmdSetEnergy, printerShort(contrast))
	queue = append(queue, c3)

	// Set mode to image mode
	c4 := formatMessage(cmdDrawingMode, []byte{0})
	queue = append(queue, c4)

	// not entirely sure what this does
	c5 := formatMessage(cmdOtherFeedPaper, cmdImgPrintSpeed)
	queue = append(queue, c5)

	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		var bmp []byte
		var bit byte
		index := 0

		for x := b.Min.X; x < b.Max.X; x++ {
			pixel := img.At(x, y)
			r, g, b, _ := pixel.RGBA()
			if r == 0 && g == 0 && b == 0 {
				bit |= 1 << index

			} else {
				bit |= 0
			}
			index++
			if index == 8 {
				index = 0
				bmp = append(bmp, bit)
				bit = 0
			}
		}
		cc := formatMessage(cmdDrawBitmap, bmp)
		queue = append(queue, cc)
	}

	// finish the lattice, whatever that means
	c6 := formatMessage(cmdControlLattice, cmdFinishLattice)
	queue = append(queue, c6)

	// feed some empty lines
	// feed_lines = 112
	c7 := formatMessage(cmdOtherFeedPaper, cmdBlankSpeed)
	queue = append(queue, c7)

	count := 112
	for count > 0 {
		feed := min(count, 0xFF)
		queue = append(queue, formatMessage(cmdFeedPaper, printerShort(feed)))
		count = count - feed
	}

	// use a GetDevState request as a way for the printer to signal that it finished
	// printing its current job.
	c8 := formatMessage(cmdGetDevState, []byte{0x00})
	queue = append(queue, c8)

	return queue
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
