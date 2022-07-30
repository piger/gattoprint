package v2

import "image"

var (
	cmdSetQuality     uint8 = 0xA4
	cmdControlLattice uint8 = 0xA6
	cmdSetEnergy      uint8 = 0xAF
	cmdDrawingMode    uint8 = 0xBE // 1 for text, 0 for images
	cmdOtherFeedPaper uint8 = 0xBD
	cmdDrawBitmap     uint8 = 0xA2 // Line to draw. 0 bit -> don't draw pixel, 1 bit -> draw pixel

	cmdPrintLattice  []uint8 = []uint8{0xAA, 0x55, 0x17, 0x38, 0x44, 0x5F, 0x5F, 0x5F, 0x44, 0x38, 0x2C}
	cmdImgPrintSpeed []uint8 = []uint8{0x23}
	cmdFinishLattice []uint8 = []uint8{0xAA, 0x55, 0x17, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x17}
)

func formatMessage(command uint8, data []uint8) []uint8 {
	var result []uint8 = []uint8{
		0x51, 0x78, command, 0x00, uint8(len(data)), 0x00,
	}
	result = append(result, data...)
	crc := crc8(data)
	result = append(result, crc)
	result = append(result, 0xFF)
	return result
}

func printerShort(i uint16) []uint16 {
	var result []uint16 = []uint16{
		i & 0xFF, (i >> 8) & 0xFF,
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

func printImage(img image.Gray) {
	var queue [][]uint8

	// set quality to standard
	c1 := formatMessage(cmdSetQuality, []uint8{0x33})
	queue = append(queue, c1)

	// start and/or set up the lattice, whatever that is
	c2 := formatMessage(cmdControlLattice, cmdPrintLattice)
	queue = append(queue, c2)

	// var contrast int16 = 12000
	//c3 := formatMessage(cmdSetEnergy, printerShort(contrast))

	// Set mode to image mode
	c4 := formatMessage(cmdDrawingMode, []uint8{0})
	queue = append(queue, c4)

	c5 := formatMessage(cmdOtherFeedPaper, cmdImgPrintSpeed)
	queue = append(queue, c5)

	// here goes all the rows
	for i := 0; i < len(img.Pix); i += img.Stride {
		var bmp []uint8
		var bit uint8

		row := img.Pix[i : i+img.Stride]
		for _, val := range row {
			if bit%8 == 0 {
				bmp = append(bmp, 0x00)
			}

			bmp[bit/8] >>= 1

			if val == 0 {
				bmp[bit/8] |= 0x80
			} else {
				bmp[bit/8] |= 0
			}

			bit += 1
		}

		cc := formatMessage(cmdDrawBitmap, bmp)
		queue = append(queue, cc)
	}

	c6 := formatMessage(cmdControlLattice, cmdFinishLattice)
	queue = append(queue, c6)
}
