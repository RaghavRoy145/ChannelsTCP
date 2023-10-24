package main

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/raghavroy145/chan-TCP-nels/tcpchan"
)

func main() {
	receiver, err := tcpchan.NewReceiver[int](":4000")
	if err != nil {
		log.Fatal(err)
	}

	sender, err := tcpchan.NewSender[int](":3000")
	if err != nil {
		log.Fatal(err)
	}
	sender.Chan <- 100
	msg := <-receiver.Chan
	fmt.Println(msg)
}
