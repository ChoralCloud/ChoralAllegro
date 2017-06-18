package main

import (
    "testing"
    "fmt"
    "net/http"
    "encoding/json"
    "time"
    "math/rand"
    "bytes"
    "os"
)

var url = "http://localhost:3000/"

type Request struct {
    DeviceId string
    DeviceData json.RawMessage
    DeviceTimestamp int64
}

// generate a random number for the deviceId (1-100)
// generate random data for deviceData
// generate timestamp for time.now()
func genRandomData() json.RawMessage {
    deviceId := rand.Intn(100)
    deviceData := json.RawMessage(`{"data":"some random data"}`)
    timestamp := time.Now().Unix()

    req := Request{
        string(deviceId),
        deviceData,
        timestamp,
    }

    j, err := json.Marshal(req)
    if err != nil {
        fmt.Println(err)
        os.Exit(3)
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
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    resp.Body.Close()
}
