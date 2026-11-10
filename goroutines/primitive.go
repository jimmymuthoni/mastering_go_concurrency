package main

import (
	"fmt"
	"time"
)

func sayHello(message string) {
	fmt.Println(message)
}

func main() {

	go sayHello("Hello Gopher welcome to concurrency..")
	go sayHello("It is a powerful tool..")
	go sayHello("Let's understand concurrency..")

	time.Sleep(time.Second * 2)

	fmt.Println("This is the entry point of the program..")
}