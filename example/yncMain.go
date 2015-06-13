package main

import (
	"fmt"

	"github.com/lindsaymarkward/go-ync"
)

var IP string = "192.168.1.221"

func main() {

	fmt.Printf("")
		response, err := ync.Discover()
		fmt.Printf("R: %v, err: %v\n", response, err)

	//	command := "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"GET\"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>"
	//	command = "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><Main_Zone><Power_Control><Power>On</Power></Power_Control></Main_Zone></YAMAHA_AV>"
	//	var cmd string = "<?xml version\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><System><Power_Control><Power>Standby</Power></Power_Control></System></YAMAHA_AV>"
	//	command = "<YAMAHA_AV cmd=\"PUT\"><Zone_2><Power_Control><Power>On</Power></Power_Control></Zone_2></YAMAHA_AV>"
	//	command = ync.Commands["GetSERVERConfig"]
//	command := "<YAMAHA_AV cmd=\"GET\"><Zone_2><Basic_Status>GetParam</Basic_Status></Zone_2></YAMAHA_AV>"
	command := "<YAMAHA_AV cmd=\"GET\"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>"
	ip := "192.168.1.221"
	//	avr := &ync.AVR{IP: ip}
	output, status, err := ync.SendCommand(command, ip)
	fmt.Printf("o %v\ns %v\ne %v\n", output, status, err)
	//	avr.SetPower("Standby", 2)
	//	avr.ChangeVolume(2, 2)
	//	powerState, _ := avr.GetPower(2)
	//	fmt.Println(powerState)
	//	avr.TogglePower(2)
	//	muteState, _ := avr.GetMuted(1)
	//	fmt.Println(muteState)
	//	avr.ToggleMuted(2)
	//	ync.SendCommand("<YAMAHA_AV cmd=\"PUT\"><Zone_2><Volume><Mute>On/Off</Mute></Volume></Zone_2></YAMAHA_AV>", ip)
	//	avr.SetMuted("Off", 1)
}
