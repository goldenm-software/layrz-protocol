package server

func (AbPacket) isServerPacket() {}
func (AcPacket) isServerPacket() {}
func (AoPacket) isServerPacket() {}
func (ArPacket) isServerPacket() {}
func (AsPacket) isServerPacket() {}
func (AuPacket) isServerPacket() {}

type ServerPackets interface {
	isServerPacket()
	ToPacket() *string
}
