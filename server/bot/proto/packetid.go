package proto

//go:generate stringer -type=PacketId
type PacketId byte

//go:generate stringer -type=ByteCode
type ByteCode byte

const (
	InfectionInformation PacketId = iota
	AttackStart
	Command         // stop attack, start self repo, stop self repo, and many others
	CommandResponse // response
	KeepAlive
	KeepAliveResponse
)

const (
	BigEndian ByteCode = iota + 1
	LittleEndian
)
