package main

import (
	"flag"
	"fmt"
	"io"
	"sync"
	"time"

	log "github.com/s00500/env_logger"

	"github.com/google/gousb"
)

var magicStartBytes = []byte{0x52, 0x4d, 0x56, 0x54}

var useGstreamer *bool = flag.Bool("gstreamer", false, "use gstreamer")

func main() {
	flag.Parse()
	log.Info("Starting djifpvvideout")
	if useGstreamer != nil && *useGstreamer {
		log.Info("using gstreamer")
	}

	openPorts := make([]string, 0)
	openPortsMu := sync.RWMutex{}

	ctx := gousb.NewContext()
	defer ctx.Close()

	vid, pid := gousb.ID(0x2ca3), gousb.ID(0x001f)

	for {
		devs, err := ctx.OpenDevices(func(d *gousb.DeviceDesc) bool {
			// this function is called for every device present.
			// Returning true means the device should be opened.
			if d.Vendor != vid || d.Product != pid {
				return false
			}

			openPortsMu.RLock()
			defer openPortsMu.RUnlock()
			return !containsString(openPorts, fmt.Sprintf("%d.%d", d.Bus, d.Address))
		})

		log.MustFatal(err)
		if len(devs) == 0 {
			time.Sleep(time.Second * 3)
			continue
		}

		log.Info("Found ", len(devs), " devices")

		for _, devInArray := range devs {
			dev := devInArray
			openPortsMu.Lock()
			openPorts = append(openPorts, fmt.Sprintf("%d.%d", dev.Desc.Bus, dev.Desc.Address))
			openPortsMu.Unlock()

			go func() {
				log.Infof("connecting to device on %d.%d", dev.Desc.Bus, dev.Desc.Address)
				openStream(dev)
				dev.Close()
				openPortsMu.Lock()
				openPorts = deleteElement(openPorts, fmt.Sprintf("%d.%d", dev.Desc.Bus, dev.Desc.Address))
				openPortsMu.Unlock()
				log.Warnf("lost device on %d.%d", dev.Desc.Bus, dev.Desc.Address)
			}()
		}
		time.Sleep(time.Second * 3)
	}
}

func openStream(dev *gousb.Device) {
	var sink StreamSink
	if useGstreamer != nil && *useGstreamer {
		sink = GstSink{}
	} else {
		//sink = FFPlaySink{}
		sink = FifoSink{Path: fmt.Sprintf("stream%d-%d.fifo", dev.Desc.Bus, dev.Desc.Address)}
	}
	ffmpegIn, stopPlayer := sink.StartInstance()
	// claim interface
	intf, done, err := googleInterface(dev)
	if err != nil {
		log.Errorf("%s.GoogleInterface: %v", dev, err)
		return
	}

	ep, err := intf.OutEndpoint(3)
	if err != nil {
		log.Errorf("%s.OutEndpoint: %v", intf, err)
		return
	}

	inEP, err := intf.InEndpoint(4)
	if err != nil {
		log.Errorf("%s.InEndpoint: %v", intf, err)
		return
	}

	// Write data to the USB device.
	numBytes, err := ep.Write(magicStartBytes)
	if numBytes != len(magicStartBytes) {
		log.Errorf("%s.Write(%d): only %d bytes written, returned error is %v", ep, len(magicStartBytes), numBytes, err)
		return
	}
	//log.Println("bytes successfully sent to the endpoint")

	stream, err := inEP.NewStream(512, 3) // Took default form github
	if err != nil {
		log.Errorf("Could not open stream: %v", intf, err)
		return
	}

	io.Copy(ffmpegIn, stream)
	stopPlayer()
	done()
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

func containsString(all []string, one string) bool {
	for _, s := range all {
		if s == one {
			return true
		}
	}
	return false
}

func deleteElement(all []string, one string) []string {
	for index, elem := range all {
		if elem == one {
			return append(all[:index], all[index+1:]...)
		}
	}
	return all
}
