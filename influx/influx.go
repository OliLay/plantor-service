package influx

import (
	"context"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	influxwrite "github.com/influxdata/influxdb-client-go/v2/api/write"
	influxdomain "github.com/influxdata/influxdb-client-go/v2/domain"
	"log"
	"strings"
	"time"
)

const BucketName = "plantor"
const OrganizationName = "plantor_org"

func Connect(host string, port int64) influx.Client {
	return influx.NewClient(fmt.Sprintf("http://%s:%d", host, port), "")
}

func Setup(client influx.Client) {
	// TODO: fix
	//orgApi := client.OrganizationsAPI()
	//var organization *influxdomain.Organization = nil
	//var exists bool
	//if organization, exists = organizationExists(orgApi); !exists {
	//	log.Printf("Organization %s doesn't exist, creating it!", OrganizationName)
	//	organization = createOrganization(orgApi)
	//}

	//bucketsApi := client.BucketsAPI()
	//if !bucketExists(bucketsApi) {
	//	log.Printf("Bucket %s doesn't exist, creating it!", BucketName)
	//	createBucket(bucketsApi, organization)
	//}
}

func GetWriteApi(client influx.Client) influxapi.WriteAPIBlocking {
	return client.WriteAPIBlocking(OrganizationName, BucketName)
}

func CreateIntMeasurement(key string, unit string, value int) *influxwrite.Point {
	key = stringifyKey(key)
	log.Printf("Creating measurement with key %s, unit %s and value %d", key, unit, value)
	return influx.NewPointWithMeasurement(key).
		AddTag("unit", unit).
		AddField("value", value).
		SetTime(time.Now())
}

func CreateFloatMeasurement(key string, unit string, value float32) *influxwrite.Point {
	key = stringifyKey(key)
	log.Printf("Creating measurement with key %s, unit %s and value %f", key, unit, value)
	return influx.NewPointWithMeasurement(key).
		AddTag("unit", unit).
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

func organizationExists(influxOrgApi influxapi.OrganizationsAPI) (*influxdomain.Organization, bool) {
	organization, err := influxOrgApi.FindOrganizationByName(context.Background(), OrganizationName)
	return organization, err == nil
}

func bucketExists(influxBucketApi influxapi.BucketsAPI) bool {
	_, err := influxBucketApi.FindBucketByName(context.Background(), BucketName)
	return err == nil
}

func createOrganization(influxOrgApi influxapi.OrganizationsAPI) *influxdomain.Organization {
	organization, err := influxOrgApi.CreateOrganizationWithName(context.Background(), OrganizationName)
	if err != nil {
		log.Fatalf("Could not create organization '%s' on InfluxDb: %s", OrganizationName, err)
	}
	return organization
}

func createBucket(influxBucketApi influxapi.BucketsAPI, organization *influxdomain.Organization) {
	_, err := influxBucketApi.CreateBucketWithName(context.Background(), organization, BucketName)
	if err != nil {
		log.Fatalf("Could not create bucket '%s' on InfluxDb: %s", BucketName, err)
	}
}
