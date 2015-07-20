package avryamaha

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// timeout value for HTTP and UDP commands
const timeout = time.Second * 2

// min/max volumes
const MaxVolume = 16.5
const MinVolume = -80.5

type AVR struct {
	IP    string
	ID    string // serial number
	Name  string
	Model string
}

type AVRState struct {
	Power  bool
	Volume float64
	Muted  bool
	Input  string
}

// resultError is a multiple-value struct for our channel
type resultError struct {
	text []byte
	err  error
}

// Result is the top-level XML container needed
type Result struct {
	Device DeviceXML `xml:"device"`
}

// DeviceXML is the XML container that stores the data we are interested in
type DeviceXML struct {
	SerialNumber string `xml:"serialNumber"`
	ModelName    string `xml:"modelName"`
}

// ChangeVolume adjusts the volume by amount for a given zone
// amount must be a supported value +/- 0.5* 1, 2, 5
// *0.5 is done with a somewhat different command
func (r *AVR) ChangeVolume(amount float64, zone int) error {
	var command string
	zoneText := setZone(zone)
	directionText := "Up"
	if amount < 0 {
		directionText = "Down"
		amount = -1 * amount
	}
	if amount == 0.5 {
		command = fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Lvl><Val>%v</Val><Exp></Exp><Unit></Unit></Lvl></Volume></%v></YAMAHA_AV>", zoneText, directionText, zoneText)
	} else {
		command = fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Lvl><Val>%v %v dB</Val><Exp></Exp><Unit></Unit></Lvl></Volume></%v></YAMAHA_AV>", zoneText, directionText, amount, zoneText)
	}
	_, err := SendCommand(command, r.IP)
	return err
}

// SetVolume sets the volume for the given zone (takes in an int like -800 or 165)
func (r *AVR) SetVolume(volume, zone int) error {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Lvl><Val>%v</Val><Exp>1</Exp><Unit>dB</Unit></Lvl></Volume></%v></YAMAHA_AV>", zoneText, volume, zoneText)
	_, err := SendCommand(command, r.IP)
	return err
}

// SetInput sets the input for the given zone
func (r *AVR) SetInput(input string, zone int) error {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Input><Input_Sel>%v</Input_Sel></Input></%v></YAMAHA_AV>", zoneText, input, zoneText)
	_, err := SendCommand(command, r.IP)
	return err
}

// SetPower sets the power ("On" or "Standby") for a given zone
func (r *AVR) SetPower(power bool, zone int) error {
	zoneText := setZone(zone)
	powerText := "Standby"
	if power {
		powerText = "On"
	}
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Power_Control><Power>%v</Power></Power_Control></%v></YAMAHA_AV>", zoneText, powerText, zoneText)
	_, err := SendCommand(command, r.IP)
	return err
}

// TogglePower toggles the power (On/Standby) for a given zone
// toggling is actually handled by the AVR when "On/Standby" is the Power value - no need to get state
// returns the present (new) state (true/false)
func (r *AVR) TogglePower(zone int) (bool, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Power_Control><Power>On/Standby</Power></Power_Control></%v></YAMAHA_AV>", zoneText, zoneText)
	// get power state before we update, otherwise it's too slow and might be incorrect
	state, _ := r.GetPower(zone)
	_, err := SendCommand(command, r.IP)
	return !state, err
}

// ToggleMuted toggles the muted state (On/Off) for a given zone
// returns the present (new) state (true/false)
func (r *AVR) ToggleMuted(zone int) (bool, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Mute>On/Off</Mute></Volume></%v></YAMAHA_AV>", zoneText, zoneText)
	_, err := SendCommand(command, r.IP)
	state, _ := r.GetMuted(zone)
	return state, err
}

// SetMuted sets the muted state to "On" or "Off" (as passed in) for a given zone
func (r *AVR) SetMuted(state bool, zone int) error {
	zoneText := setZone(zone)
	stateText := "Off"
	if state {
		stateText = "On"
	}

	command := fmt.Sprintf("<YAMAHA_AV cmd=\"PUT\"><%v><Volume><Mute>%v</Mute></Volume></%v></YAMAHA_AV>", zoneText, stateText, zoneText)
	_, err := SendCommand(command, r.IP)
	return err
}

// GetPower returns the current state of the power (true/false) for a given zone
func (r *AVR) GetPower(zone int) (bool, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, err := SendCommand(command, r.IP)
	// could use xml.Unmarshal, but doesn't seem necessary - the YNC API shouldn't change so this is simpler
	if err != nil {
		return false, err
	}
	i := strings.Index(response, "<Power>")
	if i == -1 {
		return false, fmt.Errorf("Can't get power state")
	}
	if response[i+7:i+9] == "On" {
		return true, nil
	} else {
		return false, nil
	}
}

// GetMuted returns the current muted state (true/false) for a given zone
func (r *AVR) GetMuted(zone int) (bool, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, err := SendCommand(command, r.IP)
	i := strings.Index(response, "<Mute>")
	if i == -1 {
		return false, fmt.Errorf("Can't get muted state")
	}
	if response[i+6:i+8] == "On" {
		return true, err
	} else {
		return false, err
	}
}

// GetVolume returns the current volume for a given zone
func (r *AVR) GetVolume(zone int) (float64, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, err := SendCommand(command, r.IP)
	if err != nil {
		return 0, err
	}
	return extractVolume(response)
}

// GetInput returns the current input for a given zone
func (r *AVR) GetInput(zone int) (string, error) {
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, err := SendCommand(command, r.IP)
	if err != nil {
		return "", err
	}
	start := strings.Index(response, "<Input_Sel_Item_Info><Param>")
	end := strings.Index(response, "</Param>")
	if start == -1 || end == -1 {
		return "", fmt.Errorf("Can't get input state")
	}
	valueString := response[start+len("<Input_Sel_Item_Info><Param>") : end]
	return valueString, err
}

// GetState returns the current state of the AVR as an AVRState struct - power, volume, mute, input
func (r *AVR) GetState(zone int) (AVRState, error) {
	var state AVRState
	var err error
	zoneText := setZone(zone)
	command := fmt.Sprintf("<YAMAHA_AV cmd=\"GET\"><%v><Basic_Status>GetParam</Basic_Status></%v></YAMAHA_AV>", zoneText, zoneText)
	response, err := SendCommand(command, r.IP)
	if err != nil {
		return state, err
	}
	// get power from response string
	i := strings.Index(response, "<Power>")
	if i == -1 {
		return state, fmt.Errorf("Can't get power state")
	}
	if response[i+7:i+9] == "On" {
		state.Power = true
	} else {
		state.Power = false
	}
	// get volume from response string
	volume, err := extractVolume(response)
	if err != nil {
		return state, err
	}
	state.Volume = volume
	// get muted from response string
	i = strings.Index(response, "<Mute>")
	if i == -1 {
		return state, fmt.Errorf("Can't get muted state")
	}
	if response[i+6:i+8] == "On" {
		state.Muted = true
	} else {
		state.Muted = false
	}
	// get input from response string
	start := strings.Index(response, "<Input_Sel_Item_Info><Param>")
	end := strings.Index(response, "</Param>")
	if start == -1 || end == -1 {
		return state, fmt.Errorf("Can't get input state")
	}
	valueString := response[start+len("<Input_Sel_Item_Info><Param>") : end]
	state.Input = valueString
	return state, err
}

// GetXMLData gets the name and serial number from the Yamaha AVR's description XML file
func (avr *AVR) GetXMLData() error {
	success := make(chan resultError, 1)
	// NOTE: this URL should be right, but hopefully we can get it from the SSDP response in future
	url := "http://" + avr.IP + ":49154/MediaRenderer/desc.xml"
	var err error

	go func() {
		// request xml file at url via HTTP GET
		var response *http.Response
		response, err = http.Get(url)
		// handle GET error and non-200 status (usually 404) as errors
		if err != nil {
			log.Printf("Error with GET request: %v\n", err)
			success <- resultError{nil, err}
		} else if response.Status != "200 OK" {
			success <- resultError{nil, fmt.Errorf(response.Status)}
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			success <- resultError{body, nil}
		}
	}()

	select {
	case result := <-success:
		if result.err != nil {
			return fmt.Errorf("Error getting XML data: %v\n", result.err)
		}
		var data Result
		err = xml.Unmarshal(result.text, &data)
		if err != nil {
			return err
		}
		// save data to AVR variable
		avr.ID = data.Device.SerialNumber
		avr.Model = data.Device.ModelName

		return nil
	case <-time.After(timeout):
		return fmt.Errorf("Timeout getting data")
	}
}

// Discover uses SSDP (UDP) to find and return the IP address of the Yamaha AVR
// returns an empty string if not found ??
// TODO: right now it only finds MythTV, not the AVR
func Discover() (string, error) {
	success := make(chan string, 1)

	//	searchString := "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1900\r\n MX: 20\r\n Man: \"ssdp:discover\"\r\n ST: urn:schemas-upnp-org:device:MediaServer:1\r\n"
	searchString := "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1900\r\n MX: 20\r\n Man: \"ssdp:discover\"\r\n ST: urn:schemas-upnp-org:device:MediaRenderer:1\r\n"
	searchString = "M-SEARCH * HTTP/1.1\r\n HOST:239.255.255.250:1900\r\n MX: 10\r\n Man: \"ssdp:discover\"\r\n ST: ssdp:all\r\n "
	//	var ip string
	var response string
	var err error

	go func() {
		ssdp, _ := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
		c, _ := net.ListenPacket("udp4", ":0")
		socket := c.(*net.UDPConn)
		message := []byte(searchString)
		socket.WriteToUDP(message, ssdp)
		answerBytes := make([]byte, 1024)
		// stores result in answerBytes (pass-by-reference)
		_, _, err = socket.ReadFromUDP(answerBytes)
		if err == nil {
			response = string(answerBytes)
			fmt.Println("Found: ", response)
			//			_, _, err = socket.ReadFromUDP(answerBytes)
			//			response = string(answerBytes)
			//			fmt.Println("and", response)
			// extract IP address from full response
			//		startIndex := strings.Index(response, "LOCATION: ") + 17
			//		endIndex := strings.Index(response, "MAC: ") - 2
			//		ip = response[startIndex:endIndex]
			success <- response
		}
	}()

	select {
	case <-success:
		return response, err
	case <-time.After(timeout):
		return "", fmt.Errorf("Timeout trying to find Yamaha AVR via SSDP")
	}
}

// SendCommand sends an XML YNC command to a given ip using POST
// it returns the response, the status (200 OK or 400 Bad Request) and an error value
func SendCommand(cmd, ip string) (string, error) {
	success := make(chan resultError, 1)
	var err error
	//	var response *http.Response
	// this URL could be extracted from device description XML, <yamaha:X_controlURL>/YamahaRemoteControl/ctrl</yamaha:X_controlURL>
	url := "http://" + ip + ":80/YamahaRemoteControl/ctrl"
	// prepend command with XML tag
	cmd = "<?xml version=\"1.0\" encoding=\"utf-8\"?>" + cmd

	go func() {
		var query = []byte(cmd)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(query))
		client := &http.Client{}
		response, err := client.Do(req)
		// handle POST error and non-200 status (usually 404) as errors
		if err != nil {
			log.Printf("Command was %v\n", cmd)
			log.Printf("Error with POST request: %v\n", err)
			success <- resultError{nil, err}
		} else if response.Status != "200 OK" {
			success <- resultError{nil, fmt.Errorf(response.Status)}
		} else {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			success <- resultError{body, nil}
		}
	}()
	select {
	case response := <-success:
		if response.err != nil {
			return "", err
		}
		return string(response.text), nil
	case <-time.After(timeout):
		return "", fmt.Errorf("Timeout handling HTTP command %v\n", cmd)
	}
}

func setZone(zone int) string {
	zoneText := "Main_Zone" // default zone
	if zone > 1 {
		zoneText = "Zone_" + strconv.Itoa(zone)
	}
	return zoneText
}

// extractVolume takes a full response string and returns the volume as a float
func extractVolume(response string) (float64, error) {
	// crop response string to just volume data
	start := strings.Index(response, "<Volume>")
	end := strings.Index(response, "</Volume>")
	if start == -1 || end == -1 {
		return 0, fmt.Errorf("Can't get volume state")
	}
	response = response[start:end]
	start = strings.Index(response, "<Val>")
	end = strings.Index(response, "</Val>")
	valueString := response[start+5 : end]
	// valueString will be like "-605", convert to float like -60.5
	value, errConv := strconv.ParseFloat(valueString, 0)
	return value / 10, errConv
}

// (why? probably not helpful) could check description XML file to find services AVR provides, like
// <serviceType>urn:schemas-upnp-org:service:AVTransport:1</serviceType>
// <serviceType>urn:schemas-upnp-org:service:RenderingControl:1</serviceType>
