package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"

	"github.com/itispx/goaws/s3"
)

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"testing"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

// 	"github.com/google/uuid"
// 	"github.com/itispx/goaws/s3"
// )

func getSVC(region string) (*awss3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return awss3.NewFromConfig(cfg), nil
}

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

func TestNewSession(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	_, err := s3.NewSession(&s3.NewSessionInput{
		Region: &region,
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestNewSession_NilInput(t *testing.T) {
	t.Parallel()

	_, err := s3.NewSession(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestNewSession_EmptyInput(t *testing.T) {
	t.Parallel()

	_, err := s3.NewSession(&s3.NewSessionInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestNewSession_NilRegion(t *testing.T) {
	t.Parallel()

	_, err := s3.NewSession(&s3.NewSessionInput{
		Region: nil,
	})
	if err != nil && err.Error() != "empty input" {
		log.Println(err.Error())
		t.Error("invalid error message")
	}
}

func TestNewSession_EmptyRegion(t *testing.T) {
	t.Parallel()

	region := ""

	_, err := s3.NewSession(&s3.NewSessionInput{
		Region: &region,
	})
	if err != nil && err.Error() != "empty 'Region' param" {
		log.Println(err.Error())
		t.Error("invalid error message")
	}
}

func TestBucket_NewSession(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	bct := s3.Bucket{
		Region: &region,
	}

	_, err := bct.NewSession()
	if err != nil {
		t.Error(err.Error())
	}
}

func TestBucket_NewSessionEmptyInput(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{}

	_, err := bct.NewSession()
	if err != nil && err.Error() != "empty input" {
		log.Println(err.Error())
		t.Error("invalid error message")
	}
}

func TestBucket_NewSessionNilRegion(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{
		Region: nil,
	}

	_, err := bct.NewSession()
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_NewSessionEmptyRegion(t *testing.T) {
	t.Parallel()

	region := ""

	bct := s3.Bucket{
		Region: &region,
	}

	_, err := bct.NewSession()
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestListBuckets(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	buckets := []string{}
	for i := 0; i < 3; i++ {
		b, err := createBucket(region)
		if err != nil {
			t.Errorf("bucket %d setup fail: %s\n", i, err.Error())
		}

		buckets = append(buckets, b)
	}

	// Finish setup

	bcts, err := s3.ListBuckets(&s3.ListBucketsInput{
		Region: &region,
	})
	if err != nil {
		t.Error(err.Error())
	}

outerLoop:
	for _, b1 := range buckets {
		for _, b2 := range bcts.Buckets {
			if b1 == *b2.Name {
				continue outerLoop
			}
		}

		t.Errorf("bucket '%s' not found\n", b1)
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

func TestListBuckets_NilInput(t *testing.T) {
	t.Parallel()

	_, err := s3.ListBuckets(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestListBuckets_EmptyInput(t *testing.T) {
	t.Parallel()

	_, err := s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestListBuckets_NilRegion(t *testing.T) {
	t.Parallel()

	_, err := s3.ListBuckets(&s3.ListBucketsInput{
		Region: nil,
	})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestListBuckets_EmptyRegion(t *testing.T) {
	t.Parallel()

	region := ""

	_, err := s3.ListBuckets(&s3.ListBucketsInput{
		Region: &region,
	})
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_CreateNilInput(t *testing.T) {
	t.Parallel()

	name := uuid.New().String()
	region := "us-east-1"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
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

func TestBucket_CreateEmptyInput(t *testing.T) {
	t.Parallel()

	name := uuid.New().String()
	region := "us-east-1"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.Create(&s3.BucketCreateInput{})
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

func TestBucket_CreateNilRegion(t *testing.T) {
	t.Parallel()

	name := uuid.New().String()

	bct := s3.Bucket{
		Name:   &name,
		Region: nil,
	}

	_, err := bct.Create(&s3.BucketCreateInput{})
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_CreateNilName(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	bct := s3.Bucket{
		Name:   nil,
		Region: &region,
	}

	_, err := bct.Create(&s3.BucketCreateInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_CreateEmptyRegion(t *testing.T) {
	t.Parallel()

	name := uuid.New().String()
	region := ""

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.Create(&s3.BucketCreateInput{})
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_CreateEmptyName(t *testing.T) {
	t.Parallel()

	name := ""
	region := "us-east-1"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.Create(&s3.BucketCreateInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteNilInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err = bct.Delete(nil)
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(func() {
		if t.Failed() {
			err := deleteBucket(name, region)
			if err != nil {
				t.Errorf("cleanup fail: %s", err.Error())
			}
		}
	})
}

func TestBucket_DeleteEmptyInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err = bct.Delete(&s3.BucketDeleteInput{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Cleanup(func() {
		if t.Failed() {
			err := deleteBucket(name, region)
			if err != nil {
				t.Errorf("cleanup fail: %s", err.Error())
			}
		}
	})
}

func TestBucket_DeleteNilRegion(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: nil,
	}

	_, err := bct.Delete(&s3.BucketDeleteInput{})
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteNilName(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	bct := s3.Bucket{
		Name:   nil,
		Region: &region,
	}

	_, err := bct.Delete(&s3.BucketDeleteInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteEmptyRegion(t *testing.T) {
	t.Parallel()

	region := ""
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.Delete(&s3.BucketDeleteInput{})
	if err != nil && err.Error() != "empty 'Region' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteEmptyName(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := ""

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.Delete(&s3.BucketDeleteInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObject(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Region: &region,
		Name:   &bucket,
	}

	key := "upload-test"

	_, url, err := bct.UploadObject(&s3.BucketUploadObjectInput{
		File: &[]byte{},
		Key:  &key,
	})
	if err != nil {
		t.Error(err.Error())
	}

	assertURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)

	if url != assertURL {
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

func TestBucket_UploadObjectNilInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, _, err := bct.UploadObject(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectEmptyInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, _, err := bct.UploadObject(&s3.BucketUploadObjectInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectNilBucketName(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{
		Name: nil,
	}

	_, _, err := bct.UploadObject(nil)
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectNilFile(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	key := "key-name"

	_, _, err := bct.UploadObject(&s3.BucketUploadObjectInput{
		File: nil,
		Key:  &key,
	})
	if err != nil && err.Error() != "empty 'File' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectNilKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, _, err := bct.UploadObject(&s3.BucketUploadObjectInput{
		File: &[]byte{},
		Key:  nil,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectEmptyBucketName(t *testing.T) {
	t.Parallel()

	name := ""

	bct := s3.Bucket{
		Name: &name,
	}

	_, _, err := bct.UploadObject(nil)
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_UploadObjectEmptyKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Region: nil,
		Name:   &name,
	}

	key := ""

	_, _, err := bct.UploadObject(&s3.BucketUploadObjectInput{
		File: &[]byte{},
		Key:  &key,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObject(t *testing.T) {
	t.Parallel()

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
		Name:   &bucket,
		Region: &region,
	}

	_, err = bct.GetObject(&s3.BucketGetObjectInput{
		Key: &key,
	})
	if err != nil {
		t.Error(err.Error())
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

func TestBucket_GetObjectNilInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.GetObject(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObjectEmptyInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.GetObject(&s3.BucketGetObjectInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObjectNilName(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{
		Name: nil,
	}

	_, err := bct.GetObject(&s3.BucketGetObjectInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObjectNilKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.GetObject(&s3.BucketGetObjectInput{
		Key: nil,
	})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObjectEmptyName(t *testing.T) {
	t.Parallel()

	name := ""

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.GetObject(&s3.BucketGetObjectInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_GetObjectEmptyKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	key := ""

	_, err := bct.GetObject(&s3.BucketGetObjectInput{
		Key: &key,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObject(t *testing.T) {
	t.Parallel()

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
		Name:   &bucket,
		Region: &region,
	}

	_, err = bct.DeleteObject(&s3.BucketDeleteObjectInput{
		Key: &key,
	})
	if err != nil {
		t.Error(err.Error())
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

func TestBucket_DeleteObjectNilInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.DeleteObject(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObjectEmptyInput(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.DeleteObject(&s3.BucketDeleteObjectInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObjectNilName(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{
		Name: nil,
	}

	_, err := bct.DeleteObject(&s3.BucketDeleteObjectInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObjectNilKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.DeleteObject(&s3.BucketDeleteObjectInput{
		Key: nil,
	})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObjectEmptyName(t *testing.T) {
	t.Parallel()

	name := ""

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.DeleteObject(&s3.BucketDeleteObjectInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_DeleteObjectEmptyKey(t *testing.T) {
	t.Parallel()

	name := "bucket-name"

	bct := s3.Bucket{
		Name: &name,
	}

	key := ""

	_, err := bct.DeleteObject(&s3.BucketDeleteObjectInput{
		Key: &key,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

// TODO
// Create tests for params
func TestBucket_ListObjectsNilInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	keys := []string{}
	for i := 0; i < 3; i++ {
		key := uuid.New().String()

		err := putObject(bucket, region, key)
		if err != nil {
			t.Errorf("setup fail: %s", err.Error())
		}

		keys = append(keys, key)
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   &bucket,
		Region: &region,
	}

	out, err := bct.ListObjects(nil)
	if err != nil {
		t.Error(err.Error())
	}

outerLoop:
	for _, k := range keys {
		for _, c := range out.Contents {
			if *c.Key == k {
				continue outerLoop
			}
		}

		t.Errorf("object '%s' not found", k)
	}

	t.Cleanup(func() {
		for _, k := range keys {
			err := deleteObject(bucket, region, k)
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

func TestBucket_ListObjectsEmptyInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	keys := []string{}
	for i := 0; i < 3; i++ {
		key := uuid.New().String()

		err := putObject(bucket, region, key)
		if err != nil {
			t.Errorf("setup fail: %s", err.Error())
		}

		keys = append(keys, key)
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   &bucket,
		Region: &region,
	}

	out, err := bct.ListObjects(&s3.ListObjectsInput{})
	if err != nil {
		t.Error(err.Error())
	}

outerLoop:
	for _, k := range keys {
		for _, c := range out.Contents {
			if *c.Key == k {
				continue outerLoop
			}
		}

		t.Errorf("object '%s' not found", k)
	}

	t.Cleanup(func() {
		for _, k := range keys {
			err := deleteObject(bucket, region, k)
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

func TestBucket_ListObjectsNilName(t *testing.T) {
	t.Parallel()

	bct := s3.Bucket{
		Name: nil,
	}

	_, err := bct.ListObjects(&s3.ListObjectsInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_ListObjectsEmptyName(t *testing.T) {
	t.Parallel()

	name := ""

	bct := s3.Bucket{
		Name: &name,
	}

	_, err := bct.ListObjects(&s3.ListObjectsInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGet(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	bucket, err := createBucket(region)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	key := "test-key"
	err = putObject(bucket, region, key)
	if err != nil {
		t.Errorf("setup fail: %s", err.Error())
	}

	// Finish setup

	bct := s3.Bucket{
		Name:   &bucket,
		Region: &region,
	}

	duration := time.Hour * 1
	out, err := bct.PresignGet(&s3.PresignGetInput{
		Key:      &key,
		Duration: &duration,
	})
	if err != nil {
		t.Error(err.Error())
	}

	_, err = http.Get(out.URL)
	if err != nil {
		t.Error(err.Error())
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

func TestBucket_PresignGetNilInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.PresignGet(nil)
	if err != nil && err.Error() != "nil input" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetEmptyInput(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.PresignGet(&s3.PresignGetInput{})
	if err != nil && err.Error() != "empty input" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetNilName(t *testing.T) {
	t.Parallel()

	region := "us-east-1"

	bct := s3.Bucket{
		Name:   nil,
		Region: &region,
	}

	_, err := bct.PresignGet(&s3.PresignGetInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetNilKey(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	duration := time.Hour
	_, err := bct.PresignGet(&s3.PresignGetInput{
		Key:      nil,
		Duration: &duration,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetNilDuration(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	key := "object-key"
	_, err := bct.PresignGet(&s3.PresignGetInput{
		Key:      &key,
		Duration: nil,
	})
	if err != nil && err.Error() != "empty 'Duration' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetEmptyName(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := ""

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	_, err := bct.PresignGet(&s3.PresignGetInput{})
	if err != nil && err.Error() != "empty 'Name' param" {
		t.Error("invalid error message")
	}
}

func TestBucket_PresignGetEmptyKey(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	name := "bucket-name"

	bct := s3.Bucket{
		Name:   &name,
		Region: &region,
	}

	key := ""
	duration := time.Hour
	_, err := bct.PresignGet(&s3.PresignGetInput{
		Key:      &key,
		Duration: &duration,
	})
	if err != nil && err.Error() != "empty 'Key' param" {
		t.Error("invalid error message")
	}
}

// func TestBucket_PresignPut(t *testing.T) {
// 	// trust
// }
