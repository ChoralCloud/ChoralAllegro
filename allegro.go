package main

import (
    "fmt"
    "time"
    "log"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "github.com/julienschmidt/httprouter"
)

type Data struct {
    DeviceId        string              `json:"device_id"`
    DeviceData      json.RawMessage     `json:"device_data"`
    DeviceTimestamp int64               `json:"device_timestamp"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "Choral Device Endpoint\n")
}

func checkTimestamp(timestamp int64) bool {
    curTime := time.Now().Unix()
    diff := curTime - timestamp
    fmt.Println(diff)
    if diff < 0 || diff > (60*60) {
        return false
    }
    return true
}

func checkId() bool {
    return true
}

func checkData() bool {
    return true
}

// Recieves data, in format of {deviceId, data, timestamp}
// We need to verify the format of these:
//   - Is deviceId correct?
//   - Is timestamp within the last x seconds?
//   - Is data format?
// then we need to pass it to kafka, which is in another docker container right now
func VerifyPayloadAndSend(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    log.Printf("Post requests work!")

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    payload := Data{}
    json.Unmarshal(body, &payload)

    // do checks concurrently with go func()?
    if !checkTimestamp(payload.DeviceTimestamp) || !checkId() || !checkData() {
        fmt.Fprintf(w, "Data has incorrect format")
    }

    fmt.Println(payload.DeviceId)
    fmt.Println(string(payload.DeviceData))
    fmt.Println(payload.DeviceTimestamp)
}

// Basic handlers to deal with different routes.
// All requests should come into the same route as POST requests
func handleRequests() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/", VerifyPayloadAndSend)
    log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
    handleRequests()
}
