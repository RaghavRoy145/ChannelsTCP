package main

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/raghavroy145/chan-TCP-nels/tcpchan"
)

func main() {
	receiver, err := tcpchan.NewReceiver[int](":3000")
	if err != nil {
		log.Fatal(err)
	}

	sender, err := tcpchan.NewSender[int](":3000")
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 2)
	sender.Chan <- 100
	msg := <-receiver.Chan
	fmt.Println("received from channel over the wire: ", msg)
}
