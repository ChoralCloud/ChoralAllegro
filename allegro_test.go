package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

var url = "http://localhost:3000/"

type Request struct {
	DeviceId        string          `json:"device_id"`
	UserSecret      string          `json:"user_secret"`
	DeviceData      json.RawMessage `json:"device_data"`
	DeviceTimestamp int64           `json:"device_timestamp"`
}

// generates a random UUID (this will actually come from the device itself)
func genRandomId() string {
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// generate random string for device id
// put random data in
// generate timestamp for time.now()
func genRandomData() json.RawMessage {
	rand.Seed(time.Now().UTC().UnixNano())
	deviceId := genRandomId()
	userSecret := "some_secret"
	deviceData := json.RawMessage(`{"data":"some random data"}`)
	timestamp := time.Now().Unix()

	r := Request{
		deviceId,
		userSecret,
		deviceData,
		timestamp,
	}

	j, err := json.Marshal(&r)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return j
}

// Test sending JSON to running server
// **We can use httptest to actually check the response from the server
// This would allow us to send error codes back if verification failed,
// So we could write actual test cases... but we can't do that if we want
// to use httprouter, which really simplifies the http request handling
func TestJSONFormat(*testing.T) {
	fmt.Println("Testing JSON format")
	JSON := genRandomData()
	time.Sleep(3000 * time.Millisecond)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(JSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Do(req)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
