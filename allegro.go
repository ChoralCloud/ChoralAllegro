package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/julienschmidt/httprouter"
    "encoding/json"
    "io/ioutil"
)

type Data struct {
    DeviceId        string              `json:"device_id"`
    DeviceData      json.RawMessage     `json:"device_data"`
    DeviceTimestamp int64               `json:"device_timestamp"`
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprintf(w, "Choral Device Endpoint\n")
}

func VerifyPayloadAndSend(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    log.Printf("Post requests work!")

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    payload := Data{}
    json.Unmarshal(body, &payload)
    fmt.Println(payload.DeviceId)
    fmt.Println(string(payload.DeviceData))
    fmt.Println(payload.DeviceTimestamp)
}

func handleRequests() {
    router := httprouter.New()
    router.GET("/", Index)
    router.POST("/", VerifyPayloadAndSend)
    log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
    handleRequests()
}
