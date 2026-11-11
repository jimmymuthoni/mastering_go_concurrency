// this concurrency pattern invloves looping over slice/array and adding items to the channel
package main

import "fmt"


func main(){

	myChannel := make(chan string, 3)
	names := []string{"Jim","Jay","Jason"}

	for _, name := range names {
		select {
		case myChannel <- name:
		}
		// myChannel <- name
	}
	close(myChannel) // channel is always closed by sender

	for result := range myChannel{ // we can loop over the closed channel
		fmt.Println(result)
	}

}