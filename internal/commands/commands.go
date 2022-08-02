package commands

import (
	"bytes"
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

	cmdPrintLattice  = []byte{0xAA, 0x55, 0x17, 0x38, 0x44, 0x5F, 0x5F, 0x5F, 0x44, 0x38, 0x2C}
	cmdImgPrintSpeed = []byte{0x23}
	cmdFinishLattice = []byte{0xAA, 0x55, 0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17}
	cmdBlankSpeed    = []byte{0x19}

	qualityStandard = []byte{0x33}

	cmdStart = []byte{0x51, 0x78}
	cmdEnd   = []byte{0xFF}

	// QUALITIES
	// 0x31 0x32 0x33 0x34 0x35 (49, 50, 51, 52, 53 in Java)

	// PRINT_COMMANDS
	// 0xA2 0xBF (fixed length, run-length)
)

func formatMessage(cmd byte, data []byte) []byte {
	msg := new(bytes.Buffer)

	msg.Write(cmdStart)
	msg.WriteByte(cmd)
	msg.WriteByte(0x00)
	msg.WriteByte(byte(len(data)))
	msg.WriteByte(0x00)
	msg.Write(data)
	msg.WriteByte(crc8(data))
	msg.Write(cmdEnd)

	return msg.Bytes()
}

// binary.Write(buf, binary.LittleEndian, i)
func printerShort(i int) []byte {
	result := []byte{
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

// encodeImgRows encodes each row of an image as an array of bytes; pixels
// are stored
func encodeImgRows(img *image.Gray) chan []byte {
	out := make(chan []byte)

	go func() {
		bounds := img.Bounds()

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			row := new(bytes.Buffer)
			var pixels byte
			index := 0

			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()

				if r == 0 && g == 0 && b == 0 {
					pixels |= 1 << index
				} else {
					pixels |= 0
				}

				index++

				if index == 8 {
					row.WriteByte(pixels)
					index = 0
					pixels = 0
				}
			}
			out <- row.Bytes()
		}
		close(out)
	}()
	return out
}

func PrintImage(img *image.Gray) [][]byte {
	var queue [][]byte

	// set quality to standard
	queue = append(queue, formatMessage(cmdSetQuality, qualityStandard))

	// start and/or set up the lattice, whatever that is
	queue = append(queue, formatMessage(cmdControlLattice, cmdPrintLattice))

	// Set energy used
	var contrast int = 12000
	queue = append(queue, formatMessage(cmdSetEnergy, printerShort(contrast)))

	// Set mode to image mode
	queue = append(queue, formatMessage(cmdDrawingMode, []byte{0}))

	// not entirely sure what this does
	queue = append(queue, formatMessage(cmdOtherFeedPaper, cmdImgPrintSpeed))

	// encode image, one row at a time
	for row := range encodeImgRows(img) {
		queue = append(queue, formatMessage(cmdDrawBitmap, row))
	}

	// finish the lattice, whatever that means
	queue = append(queue, formatMessage(cmdControlLattice, cmdFinishLattice))

	// feed some empty lines
	// feed_lines = 112
	queue = append(queue, formatMessage(cmdOtherFeedPaper, cmdBlankSpeed))

	count := 112
	for count > 0 {
		feed := min(count, 0xFF)
		queue = append(queue, formatMessage(cmdFeedPaper, printerShort(feed)))
		count -= feed
	}

	// use a GetDevState request as a way for the printer to signal that it finished
	// printing its current job.
	queue = append(queue, formatMessage(cmdGetDevState, []byte{0x00}))

	return queue
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
