package client

func (PaPacket) isClientPacket() {}
func (PbPacket) isClientPacket() {}
func (PcPacket) isClientPacket() {}
func (PdPacket) isClientPacket() {}
func (PiPacket) isClientPacket() {}
func (PmPacket) isClientPacket() {}
func (PrPacket) isClientPacket() {}
func (PsPacket) isClientPacket() {}

type ClientPackets interface {
	isClientPacket()
	ToPacket() *string
}
