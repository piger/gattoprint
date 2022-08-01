package bt

import "tinygo.org/x/bluetooth"

type Printer struct {
	Address bluetooth.Addresser
	device  *bluetooth.Device
}

type PrinterConfig struct {
	ServiceUUID bluetooth.UUID
	PrintUUID   bluetooth.UUID
	StatusUUID  bluetooth.UUID
}

type PrinterCommand []byte

func (p *Printer) getPrintService(uuid bluetooth.UUID) (*bluetooth.DeviceService, error) {
	svcs, err := p.device.DiscoverServices([]bluetooth.UUID{uuid})
	if err != nil {
		return nil, err
	}

	return &svcs[0], nil
}

func (p *Printer) getWriteCharacteristic(service *bluetooth.DeviceService, uuid bluetooth.UUID) (*bluetooth.DeviceCharacteristic, error) {
	chs, err := service.DiscoverCharacteristics([]bluetooth.UUID{uuid})
	if err != nil {
		return nil, err
	}
	return &chs[0], nil
}
