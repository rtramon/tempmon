package main

func TempMessage(value int) []byte {
	ipcmsg := make([]byte, 0, 32)
	ipcmsg = append(ipcmsg, EncodeByte(9)...)  // 9 = SensorDataMess
	ipcmsg = append(ipcmsg, EncodeByte(0)...)  // 0 is version
	ipcmsg = append(ipcmsg, EncodeShort(1)...) // flag TEMPERATURE = 1
	ipcmsg = append(ipcmsg, EncodeDouble(float64(value)/16.0)...)
	return ipcmsg
}

func DetectMessage(value int) []byte {
	ipcmsg := make([]byte, 0, 32)
	ipcmsg = append(ipcmsg, EncodeByte(9)...)
	ipcmsg = append(ipcmsg, EncodeByte(0)...)
	ipcmsg = append(ipcmsg, EncodeShort(8)...) // flag DETECTION = 8
	ipcmsg = append(ipcmsg, EncodeInt(value)...)
	return ipcmsg
}

func BatteryMessage(value int) []byte {
	ipcmsg := make([]byte, 0, 32)
	ipcmsg = append(ipcmsg, EncodeByte(9)...)   // 9 = SensorDataMessage
	ipcmsg = append(ipcmsg, EncodeByte(0)...)   // 0 is version
	ipcmsg = append(ipcmsg, EncodeShort(64)...) // flag BATTERY = 64
	ipcmsg = append(ipcmsg, EncodeDouble(5.0*float64(value)/1024.0)...)
	return ipcmsg
}

func IrdataMessage(value int) []byte {
	ipcmsg := make([]byte, 0, 32)
	ipcmsg = append(ipcmsg, EncodeByte(9)...)  // 9 = SensorDataMessage
	ipcmsg = append(ipcmsg, EncodeByte(0)...)  // 0 is version
	ipcmsg = append(ipcmsg, EncodeShort(4)...) // flag INTENSITY = 4
	ipcmsg = append(ipcmsg, EncodeDouble(float64(value)/10.24)...)
	return ipcmsg
}

func InitMessage(addr string) []byte {
	ipcmsg := make([]byte, 0, 32)
	ipcmsg = append(ipcmsg, EncodeByte(3)...)      // 3 = SensorMessage
	ipcmsg = append(ipcmsg, EncodeByte(0)...)      // 0 is version
	ipcmsg = append(ipcmsg, EncodeString(addr)...) // addr
	ipcmsg = append(ipcmsg, EncodeString("")...)   // name
	ipcmsg = append(ipcmsg, EncodeShort(0)...)     // deviceid
	ipcmsg = append(ipcmsg, EncodeShort(0)...)     // capabilities
	ipcmsg = append(ipcmsg, EncodeShort(0)...)     // version
	return ipcmsg
}
