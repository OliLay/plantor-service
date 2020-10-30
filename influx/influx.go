package influx

import (
	"context"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxapi "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/domain"
	"log"
	"time"
)

const BucketName = "plantor"
const OrganizationName = "plantor_org"

func Connect(host string, port int64) influx.Client {
	return influx.NewClient(fmt.Sprintf("http://%s:%d", host, port), "")
}

func Setup(client influx.Client) {
	orgApi := client.OrganizationsAPI()
	var organization *domain.Organization = nil
	var exists bool
	if organization, exists = organizationExists(orgApi); !exists {
		organization = createOrganization(orgApi)
	}

	bucketsApi := client.BucketsAPI()
	if !bucketExists(bucketsApi) {
		createBucket(bucketsApi, organization)
	}
}

func GetWriteApi(client influx.Client) influxapi.WriteAPIBlocking {
	return client.WriteAPIBlocking(OrganizationName, BucketName)
}

func PersistMeasurement(influxWriteApi influxapi.WriteAPIBlocking, key string, value float32) {
	point := influx.NewPointWithMeasurement(key).
		AddTag("unit", key).
		AddField("value", value).
		SetTime(time.Now())

	if err := influxWriteApi.WritePoint(context.Background(), point); err == nil {
		log.Printf("Wrote measurement '%s' with value '%f'", key, value)
	} else {
		log.Print("Could not write point to InfluxDb: ", err)
	}
}

func organizationExists(influxOrgApi influxapi.OrganizationsAPI) (*domain.Organization, bool) {
	organization, err := influxOrgApi.FindOrganizationByName(context.Background(), OrganizationName)
	return organization, err != nil
}

func bucketExists(influxBucketApi influxapi.BucketsAPI) bool {
	_, err := influxBucketApi.FindBucketByName(context.Background(), BucketName)
	return err != nil
}

func createOrganization(influxOrgApi influxapi.OrganizationsAPI) *domain.Organization {
	organization, err := influxOrgApi.CreateOrganizationWithName(context.Background(), OrganizationName)
	if err != nil {
		log.Fatalf("Could not create organization '%s' on InfluxDb: %s", OrganizationName, err)
	}
	return organization
}

func createBucket(influxBucketApi influxapi.BucketsAPI, organization *domain.Organization) {
	_, err := influxBucketApi.CreateBucketWithName(context.Background(), organization, BucketName)
	if err != nil {
		log.Fatalf("Could not create bucket '%s' on InfluxDb: %s", BucketName, err)
	}
}
