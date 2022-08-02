# gattoprint

A Go client for the thermal printer knows as "cat printer".

This project implements the "cat printer" protocol in Go and is based on the amazing work
done by the authors of the projects linked below.

Bluetooth support is brought by the [Go Bluetooth](https://github.com/tinygo-org/bluetooth) library,
which is part of the [TinyGo](https://tinygo.org/) project.

Installation:

```
$ go install github.com/piger/gattoprint@latest
```

## A word of warning

This project is based on the reverse engineering of the Android app iPrint, which is
the official client for this printer; using this program instead of the official client might
damage your printer.

## OS support

Tested on:

- Debian 11 (Raspbian GNU/Linux 11 (bullseye)) on a Raspberry Pi 4b.
- macOS Monterey (12.5)

Due to the requirements of the [Go Bluetooth](https://github.com/tinygo-org/bluetooth) library
cross-compiling this program is not supported, or at least is not as straightforward as it's usually
for Go programs.

## Additional resources

- [WerWolv/PythonCatPrinter](https://github.com/WerWolv/PythonCatPrinter) -- the original project that seem to have [started it all](https://werwolv.net/blog/cat_printer).
- [amber-sixel/gb01print](https://github.com/amber-sixel/gb01print) -- a fork of the above project.
- [catprinter](https://github.com/rbaron/catprinter) -- a Python client for the printer.
- [Thermal Printer wiki](https://github.com/bitbank2/Thermal_Printer/wiki/Cat-Rabbit-printer-info) -- some information about the printing protocol.
- [iPrint Utility](https://mywk.net/software/iprint-utility) -- a Windows client for the printer.
