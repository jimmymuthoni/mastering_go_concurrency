//Select allows go routines to wait on multiple communication operations.

package main

import "fmt"

func main() {

	mychannel := make(chan string)
	anotherchannel := make(chan string)

	go func () {
		mychannel <- "This data is written to channel one..."

	}()

	go func() {
		anotherchannel <- "data written to channel two...."
	}()
	
	select {
		
	case messagefrommychannel := <- mychannel:
		fmt.Println(messagefrommychannel)

	case messagefromanotherchannel := <-anotherchannel:
		fmt.Println(messagefromanotherchannel)
		
	}
}