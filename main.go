package main

import (
	"fmt"
	"io"
	"time"

	log "github.com/s00500/env_logger"

	"github.com/google/gousb"
)

var magicStartBytes = []byte{0x52, 0x4d, 0x56, 0x54}

func main() {
	log.Info("Starting")

	// TODO:
	// Open video stream
	// Write videostream to fifo, alternatively start omx via dbus and instantly output!

	ctx := gousb.NewContext()
	defer ctx.Close()

	ffmpegIn := StartFFMPEGInstance()

	for {
		// Open any device with a given VID/PID using a convenience function.
		dev, err := ctx.OpenDeviceWithVIDPID(0x2ca3, 0x001f)
		if err != nil {
			log.Errorf("Could not open googles, retrying in 2 seconds: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}
		if dev == nil {
			log.Error("No device found, retrying in 2 seconds")
			time.Sleep(time.Second * 2)
			continue
		}

		// claim interface
		intf, done, err := googleInterface(dev)
		if err != nil {
			log.Errorf("%s.GoogleInterface: %v", dev, err)
			continue
		}

		ep, err := intf.OutEndpoint(3)
		if err != nil {
			log.Errorf("%s.OutEndpoint: %v", intf, err)
			continue
		}

		inEP, err := intf.InEndpoint(4)
		if err != nil {
			log.Errorf("%s.InEndpoint: %v", intf, err)
			continue
		}

		// Write data to the USB device.
		numBytes, err := ep.Write(magicStartBytes)
		if numBytes != len(magicStartBytes) {
			log.Errorf("%s.Write(%d): only %d bytes written, returned error is %v", ep, len(magicStartBytes), numBytes, err)
			continue
		}
		//log.Println("bytes successfully sent to the endpoint")

		stream, err := inEP.NewStream(512, 3) // Took default form github
		if err != nil {
			log.Errorf("Could not open stream: %v", intf, err)
			continue
		}

		num, err := io.Copy(ffmpegIn, stream)
		//num, err := io.Copy(os.Stdout, stream)
		//_, err = io.Copy(ioutil.Discard, stream)
		if !log.Should(err) {
			log.Info("Wrote ", num, " bytes")
		}
		done()
		dev.Close()
	}
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
