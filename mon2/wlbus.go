package main

import (
	//"example/mon2/serial"
	"errors"
	"log"

	"github.com/tarm/serial"
)

type WlBus struct {
	port         *serial.Port
	state        State
	serialbuffer []byte
	index        int
}

const FRM_SOF uint8 = 0x7D
const FRM_ESC uint8 = 0x5F
const FRM_EOF uint8 = 0x7E

type State uint8

const (
	INIT State = iota
	FRAME
	ESCAPED
)

func (b *WlBus) Init() {
	c := &serial.Config{
		Name: "/dev/ttyUSB0",
		Baud: 57600,
		Size: 8,
	}

	var err error
	b.port, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	b.state = INIT
	b.serialbuffer = make([]byte, 128)
	b.index = 0
}

func (b *WlBus) Close() {
	b.port.Close()
}

func (b *WlBus) Recv(buf []byte) (n int, err error) {

	n, err = b.port.Read(buf)

	return n, err
}

func (b *WlBus) read() (byte, error) {
	// read a single byte from serial port buffer
	// if buffer is empty read serial port
	for b.index == len(b.serialbuffer) {
		// fill serialbuffer
		_, err := b.Recv(b.serialbuffer)
		if err != nil {
			return 0, err
		}
		b.index = 0
	}

	// return a byte
	retval := b.serialbuffer[b.index]
	b.index++
	return retval, nil
}

func (b *WlBus) GetFrame() ([]byte, error) {
	const MaxFrameLength = 80
	frame := make([]byte, MaxFrameLength)
	i := 0
	// read serial bytes until a complete frame is received
	for {
		ch, err := b.read()
		if err != nil {
			return nil, err
		}

		switch b.state {
		case INIT:
			if ch == FRM_SOF {
				b.state = FRAME
				i = 0
			}

		case FRAME:
			switch ch {
			case FRM_EOF:
				b.state = INIT

				return frame[:i], nil

			case FRM_ESC:
				b.state = ESCAPED

			case FRM_SOF:
				log.Println("Unexpected SOF in Frame, resetting")
				b.state = INIT
				i = 0

			default:
				// store the received byte
				if i < MaxFrameLength {
					frame[i] = ch
					i++
				} else {
					log.Println("frame buffer overflow, resetting")
					b.state = INIT
					i = 0
				}
			}

		case ESCAPED:
			// store the received byte
			frame[i] = ch
			i++
			b.state = FRAME
		}
	}
}

type MsgType uint8

const (
	MINIT MsgType = iota
	MIRDATA
	MTEMP
	MDETECT
	MVOLTAGE
	MCOMBINED
	MCLOCK
	MERROR
)

func checksum(frame []byte) error {
	sum := frame[1] + frame[2] + frame[3] + frame[4]
	if sum != frame[5] {
		return errors.New("checksum error")
	}
	return nil
}

func Type(frame []byte) MsgType {
	if checksum(frame) != nil {
		log.Println("message checksum error")
		return MERROR
	}
	msgtype := MsgType(frame[1] & 0x0f)

	return msgtype
}

func (t MsgType) String() string {
	switch t {
	case MINIT:
		return "Init"
	case MIRDATA:
		return "IrData"
	case MTEMP:
		return "Temp"
	case MDETECT:
		return "Detect"
	case MVOLTAGE:
		return "Volt"
	}
	return "unknown"
}

func Addr(frame []byte) uint8 {
	return uint8(frame[1] >> 4)
}

func Value(frame []byte) int {
	return int(frame[2])<<8 + int(frame[3])
}
