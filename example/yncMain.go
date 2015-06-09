package main

import (
	"fmt"

	"github.com/lindsaymarkward/go-ync"
)

var IP string = "192.168.1.221"

func main() {

//	response, err := ync.Discover()
//	fmt.Printf("R: %v, err: %v\n", response, err)

	command := "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"GET\"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>"
	//	command = "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><Main_Zone><Power_Control><Power>On</Power></Power_Control></Main_Zone></YAMAHA_AV>"
	//	var cmd string = "<?xml version\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><System><Power_Control><Power>Standby</Power></Power_Control></System></YAMAHA_AV>"

	ip := "192.168.1.221"
	output, status, err := ync.SendCommand(command, ip)
	fmt.Printf("o %v\ns %v\ne %v\n", output, status, err)
}
