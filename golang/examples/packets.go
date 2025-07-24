package main

import (
	"fmt"
	"strings"

	"github.com/goldenm-software/layrz-protocol/golang/v2"
)

func TestPackets() {
	fmt.Printf("Testing packets...\n\n")
	packet := "<Pb>1C9DC2691436;1740000984;19.4346059;-99.1802234;2240.800048828125;GENERIC;Core200S;-60;;06D0:01361469C29D1CC623020202;;6FF6;4FBD</Pb>"
	originalPacket := packet
	pb := layrzprotocol.PbPacket{}
	err := pb.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("(pb.FromPacket()).ToPacket() = %s\noriginalPacket =               %s\n", *pb.ToPacket(), originalPacket)
	fmt.Println("pb is equal? ", strings.Compare(originalPacket, *pb.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Ab>000000000000:MODEL1;000000000001:MODEL2;7DA8</Ab>"
	originalPacket = packet
	ab := layrzprotocol.AbPacket{}
	err = ab.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Println("(ab.FromPacket()).ToPacket() = ", *ab.ToPacket())
	fmt.Println("ab is equal? ", strings.Compare(originalPacket, *ab.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Ac>1;set_config;int:1234,float:12.34,bool:true,string:test;6C56;815F</Ac>"
	originalPacket = packet
	ac := layrzprotocol.AcPacket{}
	err = ac.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Println("(ac.FromPacket()).ToPacket() = ", *ac.ToPacket(), "originalPacket = ", originalPacket)
	fmt.Println("ac is equal? ", strings.Compare(originalPacket, *ac.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Pc>1739998848;1919;Cannot sniff in foreground;7DCB</Pc>"
	originalPacket = packet
	pc := layrzprotocol.PcPacket{}
	err = pc.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Println("(pc.FromPacket()).ToPacket() = ", *pc.ToPacket(), "originalPacket = ", originalPacket)
	fmt.Println("pc is equal? ", strings.Compare(originalPacket, *pc.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Ps>1739998822;configuration.distance.filter.meters:5,configuration.frequency.update.seconds:20,configuration.accuracy:best,configuration.server:development,configuration.sniff.interval:30,configuration.sniff.cooldown:30;BD6B</Ps>"
	originalPacket = packet
	ps := layrzprotocol.PsPacket{}
	err = ps.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("(ps.FromPacket()).ToPacket() = %s\noriginalPacket =               %s\n", *ps.ToPacket(), originalPacket)
	fmt.Println("ps is equal? ", strings.Compare(originalPacket, *ps.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Pd>1740081532;;;;;;;;report.code:LKSEN,fw.build:49,wifi.rssi:-61,cpu.temperature:43,io1.di:0,io2.di:0,io5.di:0,io6.di:0,io7.di:0,io14.di:0,io45.di:0,io46.di:0,io47.di:0;1ACF</Pd>"
	originalPacket = packet
	pd := layrzprotocol.PdPacket{}
	err = pd.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	fmt.Printf("(pd.FromPacket()).ToPacket() = %s\noriginalPacket =               %s\n", *pd.ToPacket(), originalPacket)
	fmt.Println("pd is equal? ", strings.Compare(originalPacket, *pd.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Pi>744DBD89B0D9;layrz.hub12.base;49;22246;1;460;0;false;2586</Pi>"
	originalPacket = packet
	pi := layrzprotocol.PiPacket{}
	err = pi.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("(pi.FromPacket()).ToPacket() = %s\noriginalPacket =               %s\n", *pi.ToPacket(), originalPacket)
	fmt.Println("pi is equal? ", strings.Compare(originalPacket, *pi.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Pa>phkenny123;;2664</Pa>"
	originalPacket = packet
	pa := layrzprotocol.PaPacket{}
	err = pa.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Println("(pa.FromPacket()).ToPacket() = ", *pa.ToPacket(), "originalPacket = ", originalPacket)
	fmt.Println("pa is equal? ", strings.Compare(originalPacket, *pa.ToPacket()) == 0)

	fmt.Println("----------------------")

	packet = "<Pr>;7F28</Pr>"
	originalPacket = packet
	pr := layrzprotocol.PrPacket{}
	err = pr.FromPacket(&packet)
	if err != nil {
		panic(err)
	}
	// fmt.Println("(pr.FromPacket()).ToPacket() = ", *pr.ToPacket(), "originalPacket = ", originalPacket)
	fmt.Println("pr is equal? ", strings.Compare(originalPacket, *pr.ToPacket()) == 0)
}
