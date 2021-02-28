package cmd

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveIPAddr("ip", "www.google.com")
	if err != nil {
		// fmt.Println("Resolve error ", Error())
		os.Exit(1)
	}
	fmt.Println("Resovle addr is ", addr.String())
}
