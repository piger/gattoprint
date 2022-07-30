package main

import (
	"image"
)

func encodeRunLengthRepetition(n, val int) []int {
	var result []int
	for n > 0x7f {
		result = append(result, 0x7f|(val<<7))
		n -= 0x7f
	}
	if n > 0 {
		result = append(result, (val<<7)|n)
	}
	return result
}

func runLengthEncode(row []int) []int {
	var result []int
	count := 0
	lastVal := -1

	for _, val := range row {
		if val == lastVal {
			count += 1
		} else {
			x := encodeRunLengthRepetition(count, lastVal)
			result = append(result, x...)
			count = 1
		}
		lastVal = val
	}

	if count > 0 {
		x := encodeRunLengthRepetition(count, lastVal)
		result = append(result, x...)
	}
	return result
}

func byteEncode(row []int) []int {
	var result []int

	bitEncode := func(chunkStart, bitIndex int) int {
		if row[chunkStart+bitIndex] == 1 {
			return 1 << bitIndex
		}
		return 0
	}

	// i = chunk_start
	// j = bit_index
	for i := 0; i < len(row); i += 8 {
		b := 0

		for j := 0; j < 8; j++ {
			b |= bitEncode(i, j)
		}
		result = append(result, b)
	}

	return result
}

// https://github.com/noiob/gb01print/blob/main/gb01print.py
func byteEncode2(row []uint8) []uint8 {
	var bmp []uint8
	var bit uint8

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
	return bmp
}

func printImage(im image.Gray, darkMode bool) {
	var mode []uint8

	if darkMode {
		mode = convert(CMD_PRINT_TEXT)
	} else {
		mode = convert(CMD_PRINT_IMG)
	}

	var data []uint8
	data = append(data, mode...)
	data = append(data, convert(CMD_GET_DEV_STATE)...)
	data = append(data, convert(CMD_SET_QUALITY_200_DPI)...)
	data = append(data, convert(CMD_LATTICE_START)...)

	for i := 0; i < len(im.Pix); i += im.Stride {
		row := im.Pix[i : i+im.Stride]
		encRow := byteEncode2(row)

		goop := []int8{81, 120, -65, 0, int8(len(encRow)), 0}
		// mmmh; int8 goes from -128 to 127!
		goop2 := []int8{0, 0xff}

		// goop3 is uint8
		goop3 := convert(goop)
		goop3 = append(goop3, encRow...)
		goop3 = append(goop3, convert(goop2)...)
		checksumGoop := checksum(goop3, 6, len(encRow))
		goop3[len(goop3)-2] = checksumGoop
	}
}
