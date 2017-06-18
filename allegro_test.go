package main

import (
    "testing"
    "fmt"
    "net/http"
    "encoding/json"
    "time"
    "math/rand"
    "bytes"
    "os/exec"
)

var url = "http://localhost:3000/"

type Request struct {
    DeviceId           string              `json:"device_id"`
    DeviceData         json.RawMessage     `json:"device_data"`
    DeviceTimestamp    int64               `json:"device_timestamp"`
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
    deviceData := json.RawMessage(`{"data":"some random data"}`)
    timestamp := time.Now().Unix()

    r := Request{
        deviceId,
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

// tests correctly formatted JSON
func TestJSONFormat(*testing.T) {
    fmt.Println("Testing JSON format")
    JSON := genRandomData()

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(JSON))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    client.Do(req)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }
}
