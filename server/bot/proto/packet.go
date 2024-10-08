package proto

type Packet struct {
	Type PacketType
	Args map[string]interface{}
}

func NewPacket(packetType PacketType, args map[string]interface{}) *Packet {
	return &Packet{
		Type: packetType,
		Args: args,
	}
}

func NewEmptyPacket(packetType PacketType) *Packet {
	return &Packet{
		Type: packetType,
		Args: make(map[string]interface{}),
	}
}
