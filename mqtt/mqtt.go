package mqtt

import (
	"encoding/binary"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"math"
	"plantor/influx"
	"time"
)

func Connect(host string, port int64) mqtt.Client {
	options := createClientOptions(host, port)
	client := mqtt.NewClient(options)
	token := client.Connect()

	for !token.WaitTimeout(3 * time.Second) {
	}

	if token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func SubscribeToTopics(mqttClient mqtt.Client, influxWriteApi influxapi.WriteAPIBlocking) {
	subscribeToTopic(mqttClient, influxWriteApi, "temperature")
	subscribeToTopic(mqttClient, influxWriteApi, "humidity")
	subscribeToTopic(mqttClient, influxWriteApi, "light")
	subscribeToTopic(mqttClient, influxWriteApi, "moisture")
}

func subscribeToTopic(mqttClient mqtt.Client, influxWriteApi influxapi.WriteAPIBlocking, key string) {
	mqttClient.Subscribe(key, 2, func(client mqtt.Client, msg mqtt.Message) {
		payloadBits := binary.LittleEndian.Uint32(msg.Payload())
		value := math.Float32frombits(payloadBits)

		influx.PersistMeasurement(influxWriteApi, key, value)
	})
}

func createClientOptions(mqttBrokerHost string, mqttBrokerPort int64) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerHost, mqttBrokerPort))
	opts.SetClientID("plantor-service")
	return opts
}
