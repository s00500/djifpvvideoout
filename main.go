package main

import (
	"fmt"
	"io"
	"os"

	"github.com/google/gousb"
	log "github.com/s00500/env_logger"
)

var magicStartBytes = []byte{0x52, 0x4d, 0x56, 0x54}

func main() {
	//log.Info("Starting")

	// TODO:
	// Open video stream
	// Write videostream to fifo, alternatively start omx via dbus and instantly output!

	// Initialize a new Context.
	ctx := gousb.NewContext()
	defer ctx.Close()

	// Open any device with a given VID/PID using a convenience function.
	dev, err := ctx.OpenDeviceWithVIDPID(0x2ca3, 0x001f)
	if err != nil {
		log.Fatalf("Could not open a googles: %v", err)
	}
	defer dev.Close()

	// Claim the default interface using a convenience function.
	// The default interface is always #0 alt #0 in the currently active
	// config.
	intf, done, err := googleInterface(dev)
	if err != nil {
		log.Fatalf("%s.GoogleInterface: %v", dev, err)
	}
	defer done()

	ep, err := intf.OutEndpoint(3)
	if err != nil {
		log.Fatalf("%s.OutEndpoint: %v", intf, err)
	}

	inEP, err := intf.InEndpoint(4)
	if err != nil {
		log.Fatalf("%s.InEndpoint: %v", intf, err)
	}

	// Write data to the USB device.
	numBytes, err := ep.Write(magicStartBytes)
	if numBytes != len(magicStartBytes) {
		log.Fatalf("%s.Write(%d): only %d bytes written, returned error is %v", ep, len(magicStartBytes), numBytes, err)
	}
	//log.Println("bytes successfully sent to the endpoint")

	stream, err := inEP.NewStream(512, 3) // Took default form github
	log.Must(err)
	num, err := io.Copy(os.Stdout, stream)
	log.MustFatal(err)
	log.Info("Wrote ", num, " bytes")
}

func googleInterface(d *gousb.Device) (intf *gousb.Interface, done func(), err error) {
	cfgNum, err := d.ActiveConfigNum()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get active config number of device %s: %v", d, err)
	}
	cfg, err := d.Config(cfgNum)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to claim config %d of device %s: %v", cfgNum, d, err)
	}
	i, err := cfg.Interface(3, 0)
	if err != nil {
		cfg.Close()
		return nil, nil, fmt.Errorf("failed to select interface #%d alternate setting %d of config %d of device %s: %v", 0, 0, cfgNum, d, err)
	}
	return i, func() {
		intf.Close()
		cfg.Close()
	}, nil
}
