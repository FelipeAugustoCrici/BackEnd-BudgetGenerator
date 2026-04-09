package storage

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var Client *minio.Client
var Bucket string

func Connect() {
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	Bucket = os.Getenv("S3_BUCKET")
	if Bucket == "" {
		Bucket = "budgetgen-uploads"
	}

	if endpoint == "" {
		log.Println("S3_ENDPOINT not set, storage disabled")
		return
	}

	// Strip https:// for minio client
	host := endpoint
	secure := false
	if len(host) > 8 && host[:8] == "https://" {
		host = host[8:]
		secure = true
	} else if len(host) > 7 && host[:7] == "http://" {
		host = host[7:]
	}

	var err error
	Client, err = minio.New(host, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Printf("warning: failed to connect to storage: %v", err)
		return
	}

	// For R2, bucket must be created manually in dashboard
	// Just verify connection by checking if bucket exists
	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, Bucket)
	if err != nil {
		log.Printf("warning: storage bucket check failed: %v", err)
		return
	}
	if !exists {
		log.Printf("warning: bucket '%s' does not exist — create it in the R2 dashboard", Bucket)
		return
	}

	log.Printf("storage connected, bucket: %s", Bucket)
}
