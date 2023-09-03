package energy

import (
	"context"
	"fmt"

	influxClientv2 "github.com/influxdata/influxdb-client-go/v2"
)

func createOrg() {
	client := influxClientv2.NewClient("http://localhost:8086", "admin")
	defer client.Close()
	// Create an organization
	org, err := client.OrganizationsAPI().CreateOrganizationWithName(context.Background(), "my_organization")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Org ID: %+v\n", org.Id)
}

func createBucket() {
	client := influxClientv2.NewClient("http://localhost:8086", "admin")
	defer client.Close()
	// Create a bucket in the organization
	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), "my_organization")
	if err != nil {
		panic(err)
	}
	bucket, err := client.BucketsAPI().CreateBucketWithName(context.Background(), org, "my_bucket")
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Bucket ID: %+v\n", bucket.Id)
	}
	// Create a bucket in the organization
	bucket, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, "my_bucket")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Bucket ID: %+v\n", bucket.Id)
}
