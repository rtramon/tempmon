package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

type MsgElement struct {
	Type  ElementType
	Bytes []byte
}

type Message []MsgElement

type ElementType int

const (
	Undefined ElementType = iota
	Integer
	String
	ShortString
	Double
	Char
	Byte
	Short
)

func Parse(msg []byte) Message {
	var i int = 0
	stack := []MsgElement{}

	for i < len(msg) {
		switch msg[i] {
		case 0x10: // Integer 32 bit
			stack = append(stack, MsgElement{Integer, msg[i+1 : i+5]})
			i += 5

		case 0x70: // Short 16 bit integer
			stack = append(stack, MsgElement{Short, msg[i+1 : i+3]})
			i += 3

		case 0x60: // Byte 8 bit
			stack = append(stack, MsgElement{Byte, msg[i+1 : i+2]})
			i += 2

		case 0x50: // Char 8 bit
			stack = append(stack, MsgElement{Char, msg[i+1 : i+2]})
			i += 2

		case 0x40: // Double 64 bit
			stack = append(stack, MsgElement{Double, msg[i+1 : i+9]})
			i += 9

		case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f:
			// short string, max 15 chars length in lower 4 bit of id
			l := int(msg[i]&0x0f) + 1
			stack = append(stack, MsgElement{String, msg[i+1 : i+l]})
			i += l

		case 0x20: // String, length upto 255 in next byte
			l := int(msg[i+1])
			stack = append(stack, MsgElement{String, msg[i+1 : i+l]})
			i += int(msg[i+1]) + 1

		case 0, 0xff: // error
			return stack

		default: // undefined
			stack = append(stack, MsgElement{Undefined, msg[i : i+1]})
			i++
		}
	}
	return stack
}

func Print(msg []MsgElement) {
	var i int = 0

	for i < len(msg) {
		switch msg[i].Type {
		case Integer:
			val := GetInt(msg, i)
			fmt.Printf("[%v] Integer: %v\n", i, val)

		case Short:
			val := GetInt(msg, i)
			fmt.Printf("[%v] Short Int: %v\n", i, val)

		case Byte:
			val := GetByte(msg, i)
			fmt.Printf("[%v] Byte: %v\n", i, val)

		case Char:
			fmt.Printf("[%v] Char", i)

		case Double:
			val := GetDouble(msg, i)
			fmt.Printf("[%v] Double: %v\n", i, val)

		case String:
			fmt.Printf("[%v] String: %v\n", i, string(msg[i].Bytes))

		case 0:
			return

		default:
			fmt.Printf("[%v] Unknown %v\n", i, msg[i])
		}
		i++
	}
}

func (msg *Message) GetType() int {
	return int(GetByte(*msg, 4))
}

func (msg *Message) GetSapi() string {
	return GetString(*msg, 2)
}

func GetInt(msg []MsgElement, i int) int {
	if i > len(msg) {
		return 0
	}

	switch msg[i].Type {
	case Integer:
		val := int(msg[i].Bytes[0]) | int(msg[i].Bytes[1])<<8 | int(msg[i].Bytes[2])<<16 | int(msg[i].Bytes[3])<<24
		return val

	case Short:
		val := uint16(msg[i].Bytes[0]) | uint16(msg[i].Bytes[1])<<8
		return int(val)
	case Byte:
		return int(msg[i].Bytes[0])
	}
	return 0
}

func EncodeInt(val int) []byte {
	ret := make([]byte, 5)
	ret[0] = byte(0x10) // INTEGER = 0x10
	ret[1] = byte(val & 0xFF)
	ret[2] = byte((val >> 8) & 0xFF)
	ret[3] = byte((val >> 16) & 0xFF)
	ret[4] = byte((val >> 24) & 0xFF)

	return ret
}

func GetByte(msg []MsgElement, i int) byte {
	if i < len(msg) {
		return msg[i].Bytes[0]
	}
	return 0
}

func EncodeByte(val int) []byte {
	ret := []byte{byte(0x60), byte(val)} // BYTE = 0x60
	return ret
}

func EncodeShort(val int) []byte {
	ret := []byte{byte(0x70), byte(val), byte(val >> 8)} // SHORT = 0x70
	return ret
}

func GetDouble(msg []MsgElement, i int) float64 {
	bits := binary.LittleEndian.Uint64(msg[i].Bytes)
	val := math.Float64frombits(bits)

	return val
}

func GetString(msg []MsgElement, i int) string {
	return string(msg[i].Bytes)
}

func EncodeString(val string) []byte {
	if len(val) < 16 {
		return EncodeShortString(val)
	}

	ret := []byte{byte(0x20), byte(len(val)), byte(len(val) >> 8)} // STRING is 0x20
	ret = append(ret, []byte(val)...)
	//fmt.Printf("%x\n", ret)
	return ret
}

func EncodeShortString(val string) []byte {
	l := len(val)
	ret := []byte{byte(0x30 + l)} // SHORTSTRING is 0x3y, where y is length of string
	ret = append(ret, []byte(val)...)
	//fmt.Printf("%x\n", ret)
	return ret
}

func EncodeDouble(val float64) []byte {
	ret := []byte{byte(0x40)}
	double := [8]byte{}
	binary.LittleEndian.PutUint64(double[:], math.Float64bits(val))
	ret = append(ret, double[:]...)
	return ret
}
