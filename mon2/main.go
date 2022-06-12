package main

import (
	"fmt"
	"log"
	"time"
)

type IpcMessage struct {
	addr string
	ssi  string
	data []byte
}

func main() {
	bus := WlBus{}
	bus.Init()

	ipcChan := make(chan IpcMessage)
	go loop(bus, ipcChan)

	// main loop
	for {
		// (re) initialize IPC
		ipc := Ipc{}
		for {
			err := ipc.Init()
			if err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		// handle ipc messages coming though the channel
		// in case of an ipc error, return to setup ipc again
		for {
			msg := <-ipcChan
			if err := ipc.Publish(msg.addr, msg.ssi, msg.data); err != nil {
				log.Println("IPC Publish fails", err)
				break
			}
		}
	}
}

//func loop(bus WlBus, ipc Ipc) {
func loop(bus WlBus, ipc chan IpcMessage) {
	for {
		msg, err := bus.GetFrame()
		if err != nil {
			log.Fatal(err)
		}

		if len(msg) != 6 || msg[0] != 0x5E {
			log.Println("Invalid frame, len:", len(msg))
			break
		}

		addr := fmt.Sprintf("wl/0%x", Addr(msg))
		t := Type(msg)
		value := Value(msg)
		fmt.Printf("** addr %s  type %s [%d]\n", addr, t, t)
		var ipcmsg []byte

		switch Type(msg) {
		case MTEMP:
			ipcmsg = TempMessage(value)

		case MVOLTAGE:
			ipcmsg = BatteryMessage(value)

		case MDETECT:
			ipcmsg = DetectMessage(value)

		case MINIT:
			ipcmsg = InitMessage(addr)

		case MIRDATA:
			ipcmsg = IrdataMessage(value)

		default:
			fmt.Println(Value(msg))
		}

		ipc <- IpcMessage{addr, "Data", ipcmsg}
		/* 		if err = ipc.Publish(addr, "Data", ipcmsg); err != nil {
			log.Println("IPC Publish fails", err)
			return
		} */

		fmt.Println()
	}
}
