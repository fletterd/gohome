package belkin

import (
	"errors"
	"fmt"
	"net"
	"strings"

	belkinExt "github.com/go-home-iot/belkin"
	"github.com/go-home-iot/connection-pool"
	"github.com/markdaws/gohome"
	"github.com/markdaws/gohome/log"
	"github.com/markdaws/gohome/zone"
)

type network struct {
	System *gohome.System
}

func (d *network) Devices(sys *gohome.System, modelNumber string) ([]*gohome.Device, error) {

	log.V("scanning belkin")
	var scanType belkinExt.DeviceType
	switch modelNumber {
	case "f7c043fc":
		scanType = belkinExt.DTMaker
	case "f7c029v2":
		scanType = belkinExt.DTInsight
	default:
		return nil, fmt.Errorf("unsupported model number: %s", modelNumber)
	}

	responses, err := belkinExt.Scan(scanType, 5)
	fmt.Printf("%+v\n", responses)
	if err != nil {
		log.V("scan err: %s", err)
		return nil, err
	}

	devices := make([]*gohome.Device, len(responses))
	for i, devInfo := range responses {
		err := devInfo.Load()

		if err != nil {
			// Keep going, try to get as many as we can
			log.V("failed to load device information: %s", err)
			continue
		}

		//fmt.Printf("%#v\n", response)
		//fmt.Printf("%#v\n", devInfo)

		dev, _ := gohome.NewDevice(
			modelNumber,
			devInfo.ModelName,
			devInfo.FirmwareVersion,
			strings.Replace(devInfo.Scan.Location, "/setup.xml", "", -1),
			"",
			devInfo.FriendlyName,
			devInfo.ModelDescription,
			nil,
			nil,
			nil,
			nil,
		)

		cmdBuilder := sys.Extensions.FindCmdBuilder(sys, dev)
		if cmdBuilder == nil {
			return nil, fmt.Errorf("unsupported command builder ID: %s", modelNumber)
		}
		dev.CmdBuilder = cmdBuilder

		z := &zone.Zone{
			Address:     "1",
			Name:        devInfo.FriendlyName,
			Description: devInfo.ModelDescription,
			DeviceID:    "",
			Type:        zone.ZTSwitch,
			Output:      zone.OTBinary,
		}
		dev.AddZone(z)

		if scanType == belkinExt.DTMaker {
			sensor := &gohome.Sensor{
				Address:     "1",
				Name:        devInfo.FriendlyName + " - sensor",
				Description: "",
				Attr: gohome.SensorAttr{
					Name:     "sensor",
					Value:    "-1",
					DataType: gohome.SDTInt,
					States: map[string]string{
						"0": "Closed",
						"1": "Open",
					},
				},
			}
			dev.AddSensor(sensor)
		}
		devices[i] = dev
	}

	return devices, nil
}

func (d *network) NewConnection(sys *gohome.System, dev *gohome.Device) (func(pool.Config) (net.Conn, error), error) {
	return nil, errors.New("unsupported method")
}
