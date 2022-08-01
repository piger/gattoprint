package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/piger/gattoprint/internal/bt"
	v2 "github.com/piger/gattoprint/v2"
	"tinygo.org/x/bluetooth"
)

var (
	flagOutput  = flag.String("output", "output.png", "Output file name, for preview")
	flagNoPrint = flag.Bool("no-print", false, "Disable printing, just create the preview")
)

// my printer: 657b44c5-d2b2-69e2-2c52-f33aecfb4a6f -70 GB03

func run(filename string) error {
	goo, err := v2.ConvertImage(filename)
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

	queue := v2.PrintImage(goo)

	if *flagNoPrint {
		return nil
	}

	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return err
	}
	addr, err := bt.FindDevice("GB03", adapter)
	if err != nil {
		return err
	}
	fmt.Println("found: ", addr)

	if err := bt.SendCommands(adapter, addr, queue); err != nil {
		return err
	}

	/*
		if err := v2.SendCommands(queue); err != nil {
			fmt.Printf("error sending commands: %s\n", err)
		}
	*/

	// NOTE: the original code "invert" the image using the "~" operator...
	// https://stackoverflow.com/questions/8305199/the-tilde-operator-in-python

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
