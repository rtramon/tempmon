package main

import (
	"fmt"
	"log"
	"net"
)

type Ipc struct {
	conn net.Conn
}

const address string = "192.168.2.201:12345"

func (ipc *Ipc) Init() error {
	var err error

	ipc.conn, err = net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ipc *Ipc) Subscribe(addr, ssi string) bool {
	fmt.Println("Subscribe")

	subscribe := []byte{0x70, 0xb, 0x0, 0x60, 0x2, 0x32, 0x2e, 0x2a, 0x32, 0x2e, 0x2a}

	_, err := ipc.conn.Write(subscribe)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func (ipc *Ipc) Publish(addr string, ssi string, data []byte) error {

	ipcmsg := ipc.IpcMessage(addr, ssi, data)

	//log.Printf("msg: %d [%x]", len(ipcmsg), ipcmsg)

	_, err := ipc.conn.Write(ipcmsg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (*Ipc) IpcMessage(addr string, ssi string, data []byte) []byte {
	msg := make([]byte, 0, 64)
	msg = append(msg, EncodeByte(1)...) // PUBLISH = 1
	msg = append(msg, EncodeString(addr)...)
	msg = append(msg, EncodeString(ssi)...)
	msg = append(msg, data...)

	ipcmsg := EncodeShort(len(msg) + 3) // 3 is size of encoded short for lenght
	ipcmsg = append(ipcmsg, msg...)
	return ipcmsg
}

func (ipc *Ipc) Receive(b []byte) int {
	cnt, err := ipc.conn.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	return cnt
}
