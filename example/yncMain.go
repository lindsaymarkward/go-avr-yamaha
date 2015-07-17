package main

import (
	"fmt"

	"github.com/lindsaymarkward/go-avr-yamaha"
)

var IP string = "192.168.1.221"

type Result struct {
	Device DeviceXML `xml:"device"`
}

type DeviceXML struct {
	FriendlyName string `xml:"friendlyName"`
	SerialNumber string `xml:"serialNumber"`
}

func main() {

	response, err := avryamaha.Discover()
	fmt.Printf(response, err)

	//	avr := &avryamaha.AVR{IP: IP}
	////	url := "http://192.168.1.222:6547/getDeviceDesc"
	//	url := "http://192.168.1.221:49154/MediaRenderer/desc.xml"
	//	if err := avr.GetData(url); err != nil {
	//		fmt.Printf("ERROR after: %v\n", err)
	//	}
	//	fmt.Println("\nAVR: ", avr)

	//	fmt.Println("Fin")
	//	file, err := os.Create("testGOpath.txt")
	//	if err != nil {
	//		// handle the error here
	//		return
	//	}
	//	defer file.Close()
	//
	//	file.WriteString("test")

	// get details from XML file (can't get # zones)

	//	result, err := ioutil.ReadFile("src/github.com/lindsaymarkward/go-ync/desc.xml")
	////	contents := string(result)
	//	if err == nil {
	//		fmt.Println(err)
	//		return err
	//	}
	////	fmt.Println(contents)
	//	var data Result
	//	err = xml.Unmarshal(result, &data)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	fmt.Printf("%#v\n", data)
	//	fmt.Printf("Name  : %v\n", data.Device.FriendlyName)
	//	fmt.Printf("Serial: %v\n", data.Device.SerialNumber)

	//	fmt.Printf("")

	//		response, err := avryamaha.Discover()
	//		fmt.Printf("R: %v, err: %v\n", response, err)

	//		command := "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"GET\"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>"
	//		command = "<?xml version=\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><Main_Zone><Power_Control><Power>On/Standby</Power></Power_Control></Main_Zone></YAMAHA_AV>"
	//	var cmd string = "<?xml version\"1.0\" encoding=\"utf-8\"?><YAMAHA_AV cmd=\"PUT\"><System><Power_Control><Power>Standby</Power></Power_Control></System></YAMAHA_AV>"
	//	command = "<YAMAHA_AV cmd=\"PUT\"><Zone_2><Power_Control><Power>On</Power></Power_Control></Zone_2></YAMAHA_AV>"
	//	command = avryamaha.Commands["GetSERVERConfig"]
	//		command := "<YAMAHA_AV cmd=\"GET\"><Zone_2><Basic_Status>GetParam</Basic_Status></Zone_2></YAMAHA_AV>"
	//	command := "<YAMAHA_AV cmd=\"GET\"><Main_Zone><Basic_Status>GetParam</Basic_Status></Main_Zone></YAMAHA_AV>"
	//	avr.SetPower("Standby", 2)
	//	avr.ChangeVolume(2, 2)
	//	powerState, _ := avr.GetPower(2)
	//	fmt.Println(powerState)
	//	avr.TogglePower(2)
	//	muteState, _ := avr.GetMuted(1)
	//	fmt.Println(muteState)
	//	avr.ToggleMuted(2)
	//	avryamaha.SendCommand("<YAMAHA_AV cmd=\"PUT\"><Zone_2><Volume><Mute>On/Off</Mute></Volume></Zone_2></YAMAHA_AV>", ip)
	//	avr.SetMuted("Off", 1)
	//	ip := "192.168.1.221"
	//		avr := &avryamaha.AVR{IP: ip}
	//	command := `<YAMAHA_AV cmd="PUT"><Zone_2><Volume><Lvl><Val>-85</Val><Exp>1</Exp><Unit>dB</Unit></Lvl></Volume></Zone_2></YAMAHA_AV>`
	//	command = `"<YAMAHA_AV cmd="GET"><NET_RADIO><List_Info>GetParam</List_Info></NET_RADIO></YAMAHA_AV>"`
	//	command = `"<YAMAHA_AV cmd="GET"><NET_RADIO><Play_Info>GetParam</Play_Info></NET_RADIO></YAMAHA_AV>"`
	//	output, err := avryamaha.SendCommand(command, ip)
	//	fmt.Printf("o %v\ne %v\n", output, err)
	//	fmt.Println(avr.SetVolume(0.99, 2))
	//	fmt.Println(avr.GetVolume(2))
	//	avr.SetInput("NET", 2)
	//	fmt.Println(avr.GetInput(2))
	//	fmt.Println(avr.GetState(2))
}
