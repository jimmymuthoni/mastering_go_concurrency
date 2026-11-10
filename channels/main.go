package main

import "fmt"

func main() {

	myChannel := make(chan string)

	go func() {
		myChannel <- "Data to be sent to the channel...."
	}()

	msg := <- myChannel
	fmt.Println(msg)
}