package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/piger/gattoprint/internal/bt"
	"github.com/piger/gattoprint/internal/commands"
	"github.com/piger/gattoprint/internal/graphics"

	"tinygo.org/x/bluetooth"
)

var (
	flagDeviceName = flag.String("printer-name", "GB03", "Name advertised by the printer")
	flagOutput     = flag.String("output", "output.png", "Output file name, for preview")
	flagNoPrint    = flag.Bool("no-print", false, "Disable printing, just create the preview")
)

func run(filename string) error {
	goo, err := graphics.ConvertImage(filename)
	if err != nil {
		return err
	}

	out, err := os.Create(*flagOutput)
	if err != nil {
		return err
	}
	defer out.Close()

	if err := png.Encode(out, goo); err != nil {
		return err
	}

	queue := commands.PrintImage(goo)

	if *flagNoPrint {
		return nil
	}

	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}
	addr, err := bt.FindDevice(*flagDeviceName, adapter)
	if err != nil {
		return err
	}
	log.Printf("found %s: %s", *flagDeviceName, addr)

	if err := bt.SendCommands(adapter, addr, queue); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("error: pass an image file\n")
		os.Exit(1)
	}

	filename := args[0]

	if err := run(filename); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
