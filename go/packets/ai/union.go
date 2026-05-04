package ai

func (ImPacket) isAiPacket() {}

type AiPackets interface {
	isAiPacket()
	ToPacket() *string
}
