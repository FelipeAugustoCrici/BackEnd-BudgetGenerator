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
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	Bucket = os.Getenv("MINIO_BUCKET")
	if Bucket == "" {
		Bucket = "uploads"
	}

	if endpoint == "" {
		log.Println("MINIO_ENDPOINT not set, storage disabled")
		return
	}

	var err error
	Client, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("warning: failed to connect to minio: %v", err)
		return
	}

	ctx := context.Background()
	exists, err := Client.BucketExists(ctx, Bucket)
	if err != nil {
		log.Printf("warning: minio bucket check failed: %v", err)
		return
	}
	if !exists {
		err = Client.MakeBucket(ctx, Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("warning: failed to create bucket: %v", err)
			return
		}
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + Bucket + `/*"]}]}`
		_ = Client.SetBucketPolicy(ctx, Bucket, policy)
	}

	log.Println("minio connected, bucket:", Bucket)
}
