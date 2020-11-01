package mqtt

import (
	"encoding/binary"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	"log"
	"plantor/influx"
	"strconv"
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
	subscribeToFloatTopic(mqttClient, influxWriteApi, "light/uv", "uv-index")
	subscribeToIntTopic(mqttClient, influxWriteApi, "light/ir", "lux")
	subscribeToIntTopic(mqttClient, influxWriteApi, "light/visible", "lux")

	subscribeToIntTopic(mqttClient, influxWriteApi, "temperature", "celsius")
	subscribeToIntTopic(mqttClient, influxWriteApi, "humidity", "dunno") // TODO
	subscribeToIntTopic(mqttClient, influxWriteApi, "moisture", "dunno") // TODO
}

func subscribeToFloatTopic(mqttClient mqtt.Client, influxWriteApi influxapi.WriteAPIBlocking, key string, unit string) {
	mqttClient.Subscribe(key, 2, func(client mqtt.Client, msg mqtt.Message) {
		if value, err := strconv.ParseFloat(string(msg.Payload()), 32); err == nil {
			measurement := influx.CreateFloatMeasurement(key, unit, float32(value))
			influx.PersistMeasurement(influxWriteApi, measurement)
		} else {
			log.Printf("Could not parse float value with payload %s", msg.Payload())
		}
	})
}

func subscribeToIntTopic(mqttClient mqtt.Client, influxWriteApi influxapi.WriteAPIBlocking, key string, unit string) {
	mqttClient.Subscribe(key, 2, func(client mqtt.Client, msg mqtt.Message) {
		value := binary.LittleEndian.Uint16(msg.Payload())

		measurement := influx.CreateIntMeasurement(key, unit, int(value))
		influx.PersistMeasurement(influxWriteApi, measurement)
	})
}

func createClientOptions(mqttBrokerHost string, mqttBrokerPort int64) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerHost, mqttBrokerPort))
	opts.SetClientID("plantor-service")
	return opts
}
