package trips

func (TePacket) isTripsPacket() {}
func (TsPacket) isTripsPacket() {}

type TripsPackets interface {
	isTripsPacket()
	ToPacket() *string
}
