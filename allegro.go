package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Data struct {
	DeviceId        string          `json:"device_id"`
	UserSecret      string          `json:"user_secret"`
	DeviceData      json.RawMessage `json:"device_data"`
	DeviceTimestamp int64           `json:"device_timestamp"`
}

var PRODUCER sarama.SyncProducer

// XXX HACK
// this is to keep track of the devices ip in case we need to ssh into the
// device for some reason, this is not for production this is only for us
// while testing the devices on the unreliable network at school
type DeviceRegistration struct {
	DeviceId string `json:"device_id"`
	DeviceIP string `json:"device_ip"`
}

var DEVICES []DeviceRegistration
var mu sync.Mutex

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Choral Device Endpoint\n")
}

func checkTimestamp(timestamp int64) bool {
	curTime := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	diff := curTime - timestamp
	fmt.Println(diff)
	if diff < 0 || diff > (1000*60*60) {
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
		fmt.Errorf("error while reading Body of incoming message: %s", err.Error())
		return
	}

	payload := Data{}

	json.Unmarshal(body, &payload)

	// do checks concurrently with go func()?
	if !checkTimestamp(payload.DeviceTimestamp) || !checkId() || !checkData() {
		fmt.Fprintf(w, "Data has incorrect format")
	}

	p, err := json.Marshal(&payload)
	if err != nil {
		fmt.Errorf("Error while remarshaling json: %s", err.Error())
		return
	}
	send(p)
}

func send(payload []byte) {
	topic := "choraldatastream"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}
	partition, offset, err := PRODUCER.SendMessage(msg)
	if err != nil {
		fmt.Errorf("error while sending message: %s", err.Error())
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	//:9092 for kafka
	//:2181 for zookeeper
}

// XXX HACK
// this is to keep track of the devices ip in case we need to ssh into the
// device for some reason, this is not for production this is only for us
// while testing the devices on the unreliable network at school
func ListDevices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	for _, device := range DEVICES {
		fmt.Fprintf(w, "%s\t\t\t%s\n", device.DeviceId, device.DeviceIP)
	}
}

// XXX HACK
// this is to keep track of the devices ip in case we need to ssh into the
// device for some reason, this is not for production this is only for us
// while testing the devices on the unreliable network at school
func UpdateDeviceList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Printf("Post requests work!")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Errorf("Error while reading Body of device message: %s", err.Error())
	}

	payload := DeviceRegistration{}
	json.Unmarshal(body, &payload)

	// if the device is already in there then just update it,
	// otherwise append it
	mu.Lock()

	defer mu.Unlock()

	for i, device := range DEVICES {
		if device.DeviceId == payload.DeviceId {
			DEVICES[i].DeviceIP = payload.DeviceIP
			return
		}
	}

	DEVICES = append(DEVICES, payload)
}

// Basic handlers to deal with different routes.
// All requests should come into the same route as POST requests
func handleRequests() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{"localhost:9092"}
	var err error
	PRODUCER, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// this is the only good time to panic
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/", VerifyPayloadAndSend)

	// XXX HACK
	// this is to keep track of the devices ip in case we need to ssh into the
	// device for some reason, this is not for production this is only for us
	// while testing the devices on the unreliable network at school
	router.GET("/device", ListDevices)
	router.POST("/device", UpdateDeviceList)
  
	log.Fatal(http.ListenAndServe(":8081", router))
}

func main() {
	handleRequests()
}
