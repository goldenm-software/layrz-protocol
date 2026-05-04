package trips

import "testing"

func TestTripsPackets_MarkerMethods(t *testing.T) {
	TePacket{}.isTripsPacket()
	TsPacket{}.isTripsPacket()
}
