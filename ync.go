package ync

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
)

// SendCommand sends an XML YNC command to a given ip using POST
// it returns the response, the status (200 OK or 400 Bad Request) and an error value
func SendCommand(cmd, ip string) (string, string, error) {
	url := "http://" + ip + ":80/YamahaRemoteControl/ctrl"
	var query = []byte(cmd)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(query))
	client := &http.Client{}
	response, err := client.Do(req)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return string(body), response.Status, err
}

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
