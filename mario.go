package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	file   *os.File
	reader *bufio.Reader
	c      chan os.Signal

	pipePath = flag.String("pipe", "mario-pipe", "path to the mario pipe")
)

func init() {
	flag.Parse()

	err := syscall.Mknod(*pipePath, syscall.S_IFIFO|0666, 0)

	if err != nil {
		log.Fatalf("error while creating named pipe: %+v\n", err)
	}

	file, err = os.OpenFile(*pipePath, os.O_RDONLY|syscall.O_NONBLOCK, os.ModeNamedPipe)

	if err != nil {
		log.Fatalf("error while opening named pipe for reading: %+v\n", err)
	}

	reader = bufio.NewReader(file)

	// set up signal channel
	c = make(chan os.Signal, 1)

	// notify us when we have to die so we can do a proper cleanup
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGSTOP)
}

func readFromPipe() {
	for {
		data, _, err := reader.ReadLine()

		if err == nil {
			fmt.Println(string(data))
		}

		time.Sleep(time.Duration(10) * time.Millisecond)
	}
}

func main() {
	// start reading from pipe
	go readFromPipe()

	// block until a signal is received.
	<-c

	// clean up
	os.Remove(*pipePath)
	os.Exit(0)
}
