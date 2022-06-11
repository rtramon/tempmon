package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	for {
		bus := WlBus{}
		bus.Init()

		ipc := Ipc{}
		for {
			err := ipc.Init()
			if err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		loop(bus, ipc)

		bus.Close()
	}
}

func loop(bus WlBus, ipc Ipc) {
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

		err = ipc.Publish(addr, "Data", ipcmsg)
		if err != nil {
			log.Println("IPC Publish fails", err)
			return
		}

		fmt.Println()
	}
}
