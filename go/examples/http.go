package main

import (
	"fmt"
	"time"

	"github.com/goldenm-software/layrz-protocol/go/v3/clients"
	"github.com/goldenm-software/layrz-protocol/go/v3/definitions"
	"github.com/goldenm-software/layrz-protocol/go/v3/packets/client"
)

func TestHttp() {
	fmt.Printf("Testing HTTP comm...\n\n")
	http := clients.HttpComm{}
	http.New(
		clients.HTTPS,
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

	pi := client.PiPacket{}
	pi.Ident = "link_test"
	pi.DeviceId = 0
	pi.FirmwareId = ""
	pi.FirmwareBuild = 0
	pi.FirmwareBranch = definitions.Stable
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

	position := definitions.Position{
		Latitude:  &latitude,
		Longitude: &longitude,
		Altitude:  &altitude,
	}

	extra := make(map[string]any)
	extra["golang"] = true
	extra["simulated"] = true
	extra["string"] = "HelloWorld"
	extra["int"] = 1234
	extra["float"] = 12.34

	pd := client.PdPacket{}
	pd.Timestamp = time.Now()
	pd.Position = &position
	pd.ExtraData = extra

	fmt.Printf("Sending %s...\n", *pd.ToPacket())
	output, err = http.Send(&pd)
	if err != nil {
		panic(err)
	}

	handleOutput(output)
}
