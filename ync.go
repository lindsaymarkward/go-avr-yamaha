package ync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type AVR struct {
	IP    string
	Zones int
}

// SendCommand sends an XML YNC command to a given ip using POST
// it returns the response, the status (200 OK or 400 Bad Request) and an error value
func SendCommand(cmd, ip string) (string, string, error) {
	url := "http://" + ip + ":80/YamahaRemoteControl/ctrl"
	// update command as valid string
	cmd = "<?xml version=\"1.0\" encoding=\"utf-8\"?>" + cmd

	var query = []byte(cmd)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return string(body), response.Status, err
}

// ChangeVolume adjusts the volume by amount for a given zone
// amount must be a supported value +/- 1, 2, 5
func (r *AVR) ChangeVolume(amount, zone int) error {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	directionText := "Up"
	if amount < 0 {
		directionText = "Down"
		amount = -1 * amount
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Lvl><Val>%v %v dB</Val><Exp></Exp><Unit></Unit></Lvl></Volume></%v></YAMAHA_AV>", zoneText, directionText, amount, zoneText)
	_, _, err := SendCommand(command, r.IP)
	return err
}

// SetPower sets the power to "On" or "Standby" (as passed in) for a given zone
func (r *AVR) SetPower(power string, zone int) error {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Power_Control><Power>%v</Power></Power_Control></%v></YAMAHA_AV>", zoneText, power, zoneText)
	_, _, err := SendCommand(command, r.IP)
	return err
}

// TogglePower toggles the power (On/Standby) for a given zone
// toggling is actually handled by the AVR when "On/Standby" is the Power value - no need to get state
func (r *AVR) TogglePower(zone int) error {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Power_Control><Power>On/Standby</Power></Power_Control></%v></YAMAHA_AV>", zoneText, zoneText)
	_, _, err := SendCommand(command, r.IP)
	return err
}

// ToggleMuted toggles the muted state (On/Off) for a given zone
func (r *AVR) ToggleMuted(zone int) error {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Mute>On/Off</Mute></Volume></%v></YAMAHA_AV>", zoneText, zoneText)
	_, _, err := SendCommand(command, r.IP)
	return err
}

// SetMuted sets the muted state to "On" or "Off" (as passed in) for a given zone
func (r *AVR) SetMuted(state string, zone int) error {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}

	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Mute>%v</Mute></Volume></%v></YAMAHA_AV>", zoneText, state, zoneText)
	_, _, err := SendCommand(command, r.IP)
	return err
}

// GetPower returns the current state of the power (true/false) for a given zone
func (r *AVR) GetPower(zone int) (bool, error) {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, _, err := SendCommand(command, r.IP)
	// could use xml.Unmarshal, but doesn't seem necessary - the YNC API shouldn't change so this is simpler
	i := strings.Index(response, "<Power>")
	if response[i+7:i+9] == "On" {
		return true, err
	} else {
		return false, err
	}
}

// GetMuted returns the current state of the power (true/false) for a given zone
func (r *AVR) GetMuted(zone int) (bool, error) {
	zoneText := "Main_Zone"
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, _, err := SendCommand(command, r.IP)
	fmt.Println(response)
	i := strings.Index(response, "<Mute>")
	if response[i+6:i+8] == "On" {
		return true, err
	} else {
		return false, err
	}
}

/*
?	ApplyVolume      func(state *channels.VolumeState) error
?	ApplyPlayURL func(url string, queue bool) error

*/

// Discover uses SSDP (UDP) to find and return the IP address of the Yamaha AVR
// returns an empty string if not found ??
// TODO: right now it only finds MythTV, not the AVR
func Discover() (string, error) {
	// TODO: Add timeout and return error when not found after a while
	searchString := "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1900\r\n MX: 10\r\n Man: \"ssdp:discover\"\r\n ST: urn:schemas-upnp-org:device:MediaRenderer:1\r\n "
	//	searchString = "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1900\r\n MX: 10\r\n Man: \"ssdp:discover\"\r\n ST: ssdp:all\r\n "
	//	var ip string
	var response string
	ssdp, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	c, _ := net.ListenPacket("udp4", ":0")
	socket := c.(*net.UDPConn)
	message := []byte(searchString)
	socket.WriteToUDP(message, ssdp)
	answerBytes := make([]byte, 1024)
	// stores result in answerBytes (pass-by-reference)
	_, _, err := socket.ReadFromUDP(answerBytes)
	if err == nil {
		response = string(answerBytes)
		// extract IP address from full response
		//		startIndex := strings.Index(response, "LOCATION: ") + 17
		//		endIndex := strings.Index(response, "MAC: ") - 2
		//		ip = response[startIndex:endIndex]
	}
	return response, err
}
