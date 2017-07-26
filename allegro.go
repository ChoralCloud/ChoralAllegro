package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Data struct {
	DeviceId        string          `json:"device_id"`
	UserSecret      string          `json:"user_secret"`
	DeviceData      json.RawMessage `json:"device_data"`
	DeviceTimestamp int64           `json:"device_timestamp"`
}

type Response struct {
	Err string `json:"error"`
}

func SendErrorResponse(err error, w http.ResponseWriter) {
	resp := Response{Err: err.Error()}
	str, err := json.Marshal(resp)

	if err != nil {
		log.Printf("Could not write marshal error response")
		return
	}
	log.Printf(string(str))
	fmt.Fprintf(w, string(str))
}

func SendSuccessResponse(w http.ResponseWriter) {
	resp := Response{}
	str, err := json.Marshal(resp)

	if err != nil {
		log.Printf("Could not write marshal success response")
		return
	}
	fmt.Fprintf(w, string(str))
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

func checkTimestamp(timestamp int64) error {
	curTime := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	diff := curTime - timestamp
	if diff < -(1000*60*60) || diff > (1000*60*60) {
		err_str := fmt.Sprintf("Invalid Timestamp: diff %v, curTime %v, timestamp %v", diff, curTime, timestamp)
		return errors.New(err_str)
	}
	return nil
}

func checkId() bool {
	return true
}

func checkData() bool {
	return true
}

func VerifyPayload(body []byte) error {
	payload := Data{}

	err := json.Unmarshal(body, &payload)
	if err != nil {
		log.Printf("Error while marshaling json: %s", err.Error())
		return err
	}

	err = checkTimestamp(payload.DeviceTimestamp)
	if err != nil {
		log.Printf(err.Error())
		return err
	}

	return nil
}

// Recieves data, in format of {deviceId, data, timestamp}
// We need to verify the format of these:
//   - Is deviceId correct?
//   - Is timestamp within the last x seconds?
//   - Is data format?
// then we need to pass it to kafka, which is in another docker container right now
func VerifyPayloadAndSend(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error while reading Body of incoming message: %s", err.Error())
		SendErrorResponse(err, w)
		return
	}

	err = VerifyPayload(body)
	if err != nil {
		SendErrorResponse(err, w)
		return
	}

	err = SendMessage(body)
	if err != nil {
		SendErrorResponse(err, w)
		return
	}

	SendSuccessResponse(w)
}

func SendMessage(payload []byte) error {
	topic := "choraldatastream"
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	}

	partition, offset, err := PRODUCER.SendMessage(msg)
	if err != nil {
		log.Printf("error while sending message: %s", err.Error())
		return err
	}

	if partition == -1 || offset == -1 {
		// we have had issues where this server would not connect to the kafka server
		// this is due to a issue with the compose file

		err = errors.New("There was an error submittion to kafka, your message was valid but not processed")
		log.Printf(err.Error())
		return err
	}

	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
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
		log.Printf("Error while reading Body of device message: %s", err.Error())
		SendErrorResponse(err, w)
		return
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

	broker := os.Getenv("KAFKA_URI")

	if broker == "" {
		broker = "localhost:9092"
	}

	brokers := []string{broker}

	fmt.Println("Connecting to ", brokers)

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

	defer func() {
		if err := PRODUCER.Close(); err != nil {
			log.Printf("Failed to close server", err)
		}
	}()

	PORT := ":8081"
	log.Printf("Starting on port %v", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}

func main() {
	handleRequests()
}
