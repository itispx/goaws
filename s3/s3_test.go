package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"
	"github.com/itispx/goaws/s3"
)

func createBucket(region string) (string, error) {
	name := uuid.New().String()

	svc, err := getSVC(region)
	if err != nil {
		return "", err
	}

	_, err = svc.CreateBucket(context.TODO(), &awss3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return "", err
	}

	return name, nil
}

func deleteBucket(name, region string) error {
	svc, err := getSVC(region)
	if err != nil {
		return err
	}

	_, err = svc.DeleteBucket(context.TODO(), &awss3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteObject(bucket, region, key string) error {
	svc, err := getSVC(region)
	if err != nil {
		return err
	}

	_, err = svc.DeleteObject(context.TODO(), &awss3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key),
	})

	return err
}

func putObject(bucket, region, key string) error {
	svc, err := getSVC(region)
	if err != nil {
		return err
	}

	b := []byte{}

	_, err = svc.PutObject(context.TODO(), &awss3.PutObjectInput{
		Body:   bytes.NewReader(b),
		Bucket: &bucket,
		Key:    aws.String(key),
	})

	return err
}

func getSVC(region string) (*awss3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return awss3.NewFromConfig(cfg), nil
}

func TestNewSession(t *testing.T) {
	_, err := s3.NewSession(&s3.NewSessionParams{
		Region: "us-east-1",
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNewSession_NilParams(t *testing.T) {
	_, err := s3.NewSession(nil)
	if err != nil && err.Error() != "missing params" {
		t.Error("invalid error message")
	}
}

func TestNewSession_EmptyParams(t *testing.T) {
	_, err := s3.NewSession(&s3.NewSessionParams{})
	if err != nil && err.Error() != "empty params" {
		t.Error("invalid error message")
	}
}

func TestBucket_NewSession(t *testing.T) {
	bct := s3.Bucket{
		Region: "us-east-1",
	}

	_, err := bct.NewSession()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestListBuckets(t *testing.T) {
	region := "us-east-1"

	buckets := []string{}
	for i := 0; i < 3; i++ {
		b, err := createBucket(region)
		if err != nil {
			t.Errorf("bucket %d setup fail: %s", i, err.Error())
		}

		buckets = append(buckets, b)
	}

	// Finish setup

	out, err := s3.ListBuckets(&s3.ListBucketsParams{
		Region: region,
	})
	if err != nil {
		t.Error(err.Error())
	}

outerLoop:
	for _, bo := range buckets {
		for _, bi := range out.Buckets {
			if *bi.Name == bo {
				continue outerLoop
			}
		}

		t.Errorf("bucket not found: %s", bo)
	}

	t.Cleanup(func() {
		for _, b := range buckets {
			err = deleteBucket(b, region)
			if err != nil {
				t.Errorf("bucket %s cleanup fail:", err.Error())
			}
		}
	})
}

func TestListBuckets_NilParams(t *testing.T) {
	region := "us-east-1"

	buckets := []string{}
	for i := 0; i < 3; i++ {
		b, err := createBucket(region)
		if err != nil {
			t.Errorf("bucket %d setup fail: %s", i, err.Error())
		}

		buckets = append(buckets, b)
	}

	// Finish setup

	_, err := s3.ListBuckets(nil)
	if err != nil && err.Error() != "missing params" {
		t.Error("invalid error message")
	}

	t.Cleanup(func() {
		for _, b := range buckets {
			err = deleteBucket(b, region)
			if err != nil {
				t.Errorf("bucket %s cleanup fail:", err.Error())
			}
		}
	})
}

func TestBucket_CreateEmptyParams(t *testing.T) {
	name := uuid.New().String()
	region := "us-east-1"

	bct := s3.Bucket{
		Name:   name,
		Region: region,
	}

	_, err := bct.Create(&s3.BucketCreateParams{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(func() {
		err := deleteBucket(name, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_CreateNilParams(t *testing.T) {
	name := uuid.New().String()
	region := "us-east-1"

	bct := s3.Bucket{
		Name:   name,
		Region: region,
	}

	_, err := bct.Create(nil)
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(func() {
		err := deleteBucket(name, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_DeleteParams(t *testing.T) {
	region := "us-east-1"
	name, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   name,
		Region: region,
	}

	_, err = bct.Delete(&s3.BucketDeleteParams{})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestBucket_DeleteNilParams(t *testing.T) {
	region := "us-east-1"
	name, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	bct := s3.Bucket{
		Name:   name,
		Region: region,
	}

	_, err = bct.Delete(nil)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestBucket_UploadObject(t *testing.T) {
	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   bucket,
		Region: region,
	}

	key := "upload-test"

	_, path, err := bct.UploadObject(&s3.BucketUploadObjectParams{
		File: &[]byte{},
		Key:  key,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	shouldBePath := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)

	if path != shouldBePath {
		t.Error("incorrect path")
	}

	t.Cleanup(func() {
		err = deleteObject(bucket, region, key)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}

		err = deleteBucket(bucket, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_GetObject(t *testing.T) {
	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	key := "get-test"

	err = putObject(bucket, region, key)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   bucket,
		Region: region,
	}

	_, err = bct.GetObject(&s3.BucketGetObjectParams{
		Key: key,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Cleanup(func() {
		err = deleteObject(bucket, region, key)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}

		err = deleteBucket(bucket, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_DeleteObject(t *testing.T) {
	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	key := "delete-test"

	err = putObject(bucket, region, key)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   bucket,
		Region: region,
	}

	_, err = bct.DeleteObject(&s3.BucketDeleteObjectParams{
		Key: key,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Cleanup(func() {
		err = deleteBucket(bucket, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_ListObjects(t *testing.T) {
	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	ids := []string{}
	for i := 0; i < 3; i++ {
		id := uuid.New()

		err := putObject(bucket, region, id.String())
		if err != nil {
			t.Errorf("setup fail: %s", err.Error())
		}

		ids = append(ids, id.String())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   bucket,
		Region: region,
	}

	out, err := bct.ListObjects(&s3.ListObjectsParams{})
	if err != nil {
		t.Errorf(err.Error())
	}

idsLoop:
	for _, i := range ids {
		for _, c := range out.Contents {
			if *c.Key == i {
				continue idsLoop
			}
		}

		t.Errorf("failed to list object: %s", i)
	}

	t.Cleanup(func() {
		for _, id := range ids {
			err := deleteObject(bucket, region, id)
			if err != nil {
				t.Errorf("cleanup fail: %s", err.Error())
			}
		}

		err = deleteBucket(bucket, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_PresignGet(t *testing.T) {
	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	key := "get-test"

	err = putObject(bucket, region, key)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   bucket,
		Region: region,
	}

	out, err := bct.PresignGet(&s3.PresignGetParams{
		Key:      key,
		Duration: time.Minute * 1,
	})
	if err != nil {
		t.Errorf(err.Error())
	}

	_, err = http.Get(out.URL)
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Cleanup(func() {
		err = deleteObject(bucket, region, key)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}

		err = deleteBucket(bucket, region)
		if err != nil {
			t.Errorf("cleanup fail: %s", err.Error())
		}
	})
}

func TestBucket_PresignPut(t *testing.T) {
	// trust
}
