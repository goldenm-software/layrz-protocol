package client

import "testing"

func TestClientPackets_MarkerMethods(t *testing.T) {
	PaPacket{}.isClientPacket()
	PbPacket{}.isClientPacket()
	PcPacket{}.isClientPacket()
	PdPacket{}.isClientPacket()
	PiPacket{}.isClientPacket()
	PmPacket{}.isClientPacket()
	PrPacket{}.isClientPacket()
	PsPacket{}.isClientPacket()
}
