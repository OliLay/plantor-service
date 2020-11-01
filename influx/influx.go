package influx

import (
	"context"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	influxwrite "github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"strings"
	"time"
)

const BucketName = "plantor"
const OrganizationName = "plantor_org"

func Connect(host string, port int64) influx.Client {
	return influx.NewClient(fmt.Sprintf("http://%s:%d", host, port), "")
}

func GetWriteApi(client influx.Client) influxapi.WriteAPIBlocking {
	return client.WriteAPIBlocking(OrganizationName, BucketName)
}

func CreateIntMeasurement(key string, value int) *influxwrite.Point {
	key = stringifyKey(key)
	log.Printf("Creating measurement with key %s and value %d", key, value)
	return influx.NewPointWithMeasurement(key).
		AddField("value", value).
		SetTime(time.Now())
}

func CreateFloatMeasurement(key string, value float32) *influxwrite.Point {
	key = stringifyKey(key)
	log.Printf("Creating measurement with key %s and value %f", key, value)
	return influx.NewPointWithMeasurement(key).
		AddField("value", value).
		SetTime(time.Now())
}

func PersistMeasurement(influxWriteApi influxapi.WriteAPIBlocking, point *influxwrite.Point) {
	if err := influxWriteApi.WritePoint(context.Background(), point); err == nil {
		log.Printf("Wrote measurement of type '%s'.", point.Name())
	} else {
		log.Print("Could not write measurement to InfluxDb: ", err)
	}
}

func stringifyKey(key string) string {
	return strings.Replace(key, "/", "-", -1)
}
