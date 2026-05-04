package server

import "testing"

func TestServerPackets_MarkerMethods(t *testing.T) {
	AbPacket{}.isServerPacket()
	AcPacket{}.isServerPacket()
	AoPacket{}.isServerPacket()
	ArPacket{}.isServerPacket()
	AsPacket{}.isServerPacket()
	AuPacket{}.isServerPacket()
}
