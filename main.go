package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tarm/serial"
)

func main() {
	serialDeviceName := flag.String("dev", "/dev/serial0", "Device to open")
	bufferSize := flag.Int("bufferSize", 100, "Size of internal buffer")
	baud := flag.Int("baud", 57600, "Set the baud rate")
	flag.Parse()

	serialPortConfig := &serial.Config{Name: *serialDeviceName, Baud: *baud}
	serialDevice, err := serial.OpenPort(serialPortConfig)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		buf := make([]byte, *bufferSize)
		for {
			numBytes, _ := serialDevice.Read(buf)
			fmt.Printf("%s", buf[:numBytes])
		}
	}()

	fin := make(chan bool)

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				fin <- true
			}

			_, err = serialDevice.Write([]byte(string([]rune{char})))
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	<-fin
}
