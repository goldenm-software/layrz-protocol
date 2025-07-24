package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/goldenm-software/layrz-protocol/golang/v2"
	"github.com/matishsiao/goInfo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func TestTcp() {
	fmt.Printf("Testing Tcp comm...\n\n")
	tcp := layrzprotocol.TcpComm{}
	ident := "link_server_pruebas"
	tcp.New("<server>", 1234, ident, "")
	// tcp.New("127.0.0.1", 5000, ident, "")
	err := tcp.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %s\n", err)
	}
	tcp.SetCallback(tcpCallback)
	log.Println("Conection stablished, sending Pi packet...")

	err = tcp.Send(&layrzprotocol.PiPacket{
		Ident:          ident,
		FirmwareId:     "com.layrz.link.servers",
		FirmwareBuild:  0,
		DeviceId:       0,
		HardwareId:     0,
		ModelId:        0,
		FirmwareBranch: layrzprotocol.Stable,
		FotaEnabled:    false,
	})

	if err != nil {
		log.Fatalf("Failed to send packet: %s\n", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Stop signal received...")
				return
			default:
				pd := extractData()
				log.Printf("Sending %s\n", *pd.ToPacket())
				err = tcp.Send(&pd)
				if err != nil {
					log.Fatalf("Failed to send packet: %s\n", err)
				}
				time.Sleep(time.Second * 5)
			}
		}
	}()

	<-stop
	fmt.Println("Stopping...")
	tcp.Close()
}

func tcpCallback(data *interface{}) {
	log.Printf("Received: %v\n", data)
}

func extractData() layrzprotocol.PdPacket {
	var data map[string]interface{} = make(map[string]interface{})

	gi, err := goInfo.GetInfo()
	if err == nil {
		data["os.name"] = gi.OS
		data["os.system"] = gi.Kernel
		data["os.core"] = gi.Core
		data["os.architecture"] = gi.Platform
		data["os.hostname"] = gi.Hostname
	} else {
		data["os.system"] = runtime.GOOS
		data["os.architecture"] = runtime.GOARCH
	}

	// Get CPU information if enabled
	data["cpu.count"] = runtime.NumCPU()
	cpuInfo, err := cpu.Info()
	if err == nil {
		data["cpu.name"] = cpuInfo[0].ModelName
		data["cpu.frequency.ghrz"] = cpuInfo[0].Mhz / 1000.0
	}
	memInfo, err := mem.VirtualMemory()
	if err == nil {
		data["memory.total.gb"] = memInfo.Total / 1024 / 1024 / 1024
		data["memory.used.gb"] = memInfo.Used / 1024 / 1024 / 1024
		data["memory.free.gb"] = memInfo.Free / 1024 / 1024 / 1024
		data["memory.used.percent"] = memInfo.UsedPercent
	}

	pd := layrzprotocol.PdPacket{}
	pd.Timestamp = time.Now()
	pd.ExtraData = &data

	return pd
}
