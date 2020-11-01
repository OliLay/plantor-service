package main

import (
	"log"
	"os"
	"plantor/influx"
	"plantor/mqtt"
	"strconv"
)

const EnvMqttBrokerHost = "MQTT_BROKER_HOST"
const EnvMqttBrokerPort = "MQTT_BROKER_PORT"
const EnvInfluxHost = "INFLUX_HOST"
const EnvInfluxPort = "INFLUX_PORT"

func main() {
	mqttBrokerHost, mqttBrokerPort := getServiceHostAndPort(EnvMqttBrokerHost, EnvMqttBrokerPort)
	influxHost, influxPort := getServiceHostAndPort(EnvInfluxHost, EnvInfluxPort)

	mqttClient := mqtt.Connect(mqttBrokerHost, mqttBrokerPort)

	influxClient := influx.Connect(influxHost, influxPort)

	influxWriteApi := influx.GetWriteApi(influxClient)
	mqtt.SubscribeToTopics(mqttClient, influxWriteApi)

	select {}
}

func getServiceHostAndPort(envHostnameKey string, envPortKey string) (string, int64) {
	serviceHost, exists := os.LookupEnv(envHostnameKey)
	if !exists {
		log.Fatalf("No service host defined. Please supply it using the env variable %s.",
			envHostnameKey)
	}

	servicePortString, exists := os.LookupEnv(envPortKey)
	if !exists {
		log.Fatalf("No service port defined. Please supply it using the env variable %s.",
			envPortKey)
	}

	servicePort, err := strconv.ParseInt(servicePortString, 10, 0)

	if err == nil {
		log.Printf("Parsed service %s:%d", serviceHost, servicePort)
		return serviceHost, servicePort
	} else {
		log.Fatalf("Could not parse given service port '%s'!", servicePortString)
		return "", 0
	}
}
