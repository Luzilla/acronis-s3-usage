// Package s3 simulates an administrative interface for people who maintain ostor
// this package wraps around ostor and minio/minio-go to executes calls on behalf
// of accounts in the system. This is achieved by returning an account's credential
// pair and using it for calls. It requires that an account has one. You can call
// the key management features in the ostor CLI to achieve that or use the user methods
// in the ostor package to achieve the same.
package s3

import (
	"context"
	"fmt"
	log "log/slog"
	"net/url"
	"os"

	"github.com/Luzilla/acronis-s3-usage/pkg/ostor"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3 wraps the s3 client to provide shorthands for common ops.
type S3 struct {
	mc *minio.Client
}

// NewS3 creates an S3 handler which is specific to the provided email.
func NewS3(endpointURL, email string, ostorClient *ostor.Ostor) (*S3, error) {
	endpoint, err := url.Parse(endpointURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s: %s", endpoint, err)
	}

	user, _, err := ostorClient.GetUser(email)
	if err != nil {
		return nil, err
	}

	if len(user.AccessKeys) == 0 {
		return nil, fmt.Errorf("account has no keys, please generate a key-pair first")
	}

	keyPair := user.AccessKeys[0]

	mc, err := minio.New(endpoint.Host, &minio.Options{
		Creds: credentials.NewStaticV4(
			keyPair.AccessKeyID,
			keyPair.SecretAccessKey,
			""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize s3 client: %s", err)
	}

	return &S3{mc: mc}, nil
}

// IsDeletable determines if a bucket exists and if is empty
func (s *S3) IsDeletable(ctx context.Context, bucketName string) (status bool, err error) {
	status, err = s.mc.BucketExists(ctx, bucketName)
	if err != nil {
		err = fmt.Errorf("unable to check if the bucket %q exists: %s", bucketName, err)
		return
	}

	if !status {
		err = fmt.Errorf("bucket %q does not exist", bucketName)
		return
	}

	listChan := s.mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
		MaxKeys:   1,
	})
	_, ok := <-listChan
	if !ok {
		status = false
		err = fmt.Errorf("bucket %q is not empty", bucketName)
		return
	}

	return
}

// DeleteBucket does a recursive delete on all objects within a bucket to empty it, before deleting it.
func (s *S3) DeleteBucket(ctx context.Context, bucketName string) error {
	delChan := make(chan minio.ObjectInfo)

	go func() {
		defer close(delChan)
		for object := range s.ListContents(ctx, bucketName) {
			if object.Err != nil {
				log.Error(object.Err.Error())
				os.Exit(1)
			}
			delChan <- object
		}
	}()

	for rErr := range s.mc.RemoveObjects(ctx, bucketName, delChan, minio.RemoveObjectsOptions{}) {
		log.Error("Error detected during deletion: " + rErr.Err.Error())
	}

	return s.mc.RemoveBucket(ctx, bucketName)
}

// ListContents (recursively) lists the contents of a bucket and returns a channel to "range" on.
func (s *S3) ListContents(ctx context.Context, bucketName string) <-chan minio.ObjectInfo {
	return s.mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
}
