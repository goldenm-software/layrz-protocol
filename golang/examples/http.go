package main

import (
	"fmt"
	"time"

	"github.com/goldenm-software/layrz-protocol/golang/v2"
)

func TestHttp() {
	fmt.Printf("Testing HTTP comm...\n\n")
	http := layrzprotocol.HttpComm{}
	http.New(
		layrzprotocol.HTTPS,
		"<server>",
		"link_test",
		"",
	)

	fmt.Printf("Getting commands from the server...\n\n")
	output, err := http.GetCommands()
	if err != nil {
		panic(err)
	}

	handleOutput(output)

	pi := layrzprotocol.PiPacket{}
	pi.Ident = "link_test"
	pi.DeviceId = 0
	pi.FirmwareId = ""
	pi.FirmwareBuild = 0
	pi.FirmwareBranch = layrzprotocol.Stable
	pi.HardwareId = 0
	pi.ModelId = 0
	pi.FotaEnabled = false

	fmt.Printf("Sending %s...\n", *pi.ToPacket())
	output, err = http.Send(&pi)
	if err != nil {
		panic(err)
	}

	handleOutput(output)

	latitude := 19.4346059
	longitude := -99.1802234
	altitude := 2240.00

	position := layrzprotocol.Position{
		Latitude:  &latitude,
		Longitude: &longitude,
		Altitude:  &altitude,
	}

	extra := make(map[string]interface{})
	extra["golang"] = true
	extra["simulated"] = true
	extra["string"] = "HelloWorld"
	extra["int"] = 1234
	extra["float"] = 12.34

	pd := layrzprotocol.PdPacket{}
	pd.Timestamp = time.Now()
	pd.Position = &position
	pd.ExtraData = &extra

	fmt.Printf("Sending %s...\n", *pd.ToPacket())
	output, err = http.Send(&pd)
	if err != nil {
		panic(err)
	}

	handleOutput(output)

	// msg := "Hello world"
	// pc := layrzprotocol.PcPacket{}
	// pc.Timestamp = time.Now()
	// pc.CommandId = 1924
	// pc.Message = &msg

	// fmt.Printf("Sending %s...\n", *pc.ToPacket())
	// output, err = http.Send(&pc)
	// if err != nil {
	// 	panic(err)
	// }

	// handleHttpOutput(*output)
}
