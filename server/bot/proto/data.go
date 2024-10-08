package proto

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"net/url"
	"time"
)

//go:generate stringer -type=PacketType
type PacketType int

//go:generate stringer -type=RequestType
type RequestType int

//go:generate stringer -type=ArgType
type ArgType int32

//go:generate stringer -type=ByteCode
type ByteCode int

const (
	RequestTypeClose RequestType = iota
	RequestTypeAlreadyConnected
)

const (
	PacketTypeHandshake PacketType = iota
	PacketTypeInfo
	PacketTypeRequest
	PacketTypeKeepAlive
	PacketTypeStartAttack
	PacketTypeStopAttack
	PacketTypeExecCmd
	PacketTypeExecCmdResult
	PacketTypeCriticalErrorReport
)

const (
	ArgTypeUnknown ArgType = iota - 1
	ArgTypeInt16
	ArgTypeInt32
	ArgTypeInt64
	ArgTypeString
	ArgTypeBool
	ArgTypeBinary
	ArgTypeStringList
	ArgTypeStringMap
	ArgTypeFloat
	ArgTypeDouble
	ArgTypeBotInfo
	ArgTypeIP
	ArgTypeURL
	ArgTypeDuration
)

const (
	BigEndian ByteCode = iota + 1
	LittleEndian
)

type BotInfo struct {
	Username        string
	OS              string
	Kernel          string
	Arch            string
	Version         string
	InfectionMethod string
	Processors      int32
	UpTime          int32
	TotalMemory     float32
	FreeMemory      float32
	TimeZoneDiff    float64
	SystemTime      int64
	IsRoot          bool
	LittleEndian    bool
	Is64Bit         bool
	IsDebugMode     bool
	IsIPv6Supported bool
	HasAnySSLLib    bool
	FirewallStatus  bool
}

func (b BotInfo) String() string {
	marshal, err := json.Marshal(b)
	if err != nil {
		return err.Error()
	}

	return string(marshal)
}

func (c ByteCode) Order() binary.ByteOrder {
	switch c {
	case BigEndian:
		return binary.BigEndian
	default:
		return binary.LittleEndian
	}
}

func GetDataArgType(data interface{}) ArgType {
	switch data.(type) {
	case int16:
		return ArgTypeInt16
	case int32:
		return ArgTypeInt32
	case int64:
		return ArgTypeInt64
	case string:
		return ArgTypeString
	case bool:
		return ArgTypeBool
	case []byte:
		return ArgTypeBinary
	case []string:
		return ArgTypeStringList
	case map[string]string:
		return ArgTypeStringMap
	case float32:
		return ArgTypeFloat
	case float64:
		return ArgTypeDouble
	case BotInfo:
		return ArgTypeBotInfo
	case net.IP:
		return ArgTypeIP
	case url.URL:
		return ArgTypeURL
	case time.Duration:
		return ArgTypeDuration
	default:
		return ArgTypeUnknown
	}
}
