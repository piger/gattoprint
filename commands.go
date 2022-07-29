package main

import (
	"image"
)

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

	var rows [][]uint8
	b := im.Bounds()
	row := 0
	for i := 0; i < b.Max.Y; i += im.Stride {
		rows[row] = im.Pix[i:im.Stride]
		row++
	}
}
