package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	Name   *string `json:"name"`
	Region *string `json:"region"`
	Client *s3.Client
}

type NewSessionInput struct {
	Region *string
}

func NewSession(input *NewSessionInput) (*s3.Client, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (NewSessionInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Region == nil || *input.Region == "" {
		return nil, fmt.Errorf("empty 'Region' param")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*input.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	return svc, nil
}

func (b *Bucket) NewSession() (*s3.Client, error) {
	if *b == (Bucket{}) {
		return nil, fmt.Errorf("empty input")
	}
	if b.Region == nil || *b.Region == "" {
		return nil, fmt.Errorf("empty 'Region' param")
	}

	svc, err := NewSession(&NewSessionInput{
		Region: b.Region,
	})
	if err != nil {
		return nil, err
	}

	b.Client = svc

	return b.Client, nil
}

type ListBucketsInput struct {
	SVC    *s3.Client
	Region *string
}

func ListBuckets(input *ListBucketsInput) (*s3.ListBucketsOutput, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (ListBucketsInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Region == nil || *input.Region == "" {
		return nil, fmt.Errorf("empty 'Region' param")
	}

	if input.SVC == nil {
		var err error
		input.SVC, err = NewSession(&NewSessionInput{
			Region: input.Region,
		})
		if err != nil {
			return nil, err
		}
	}

	resp, err := input.SVC.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	return resp, nil
}

type BucketCreateInput struct {
	*s3.CreateBucketInput
}

func (b *Bucket) Create(input *BucketCreateInput) (*s3.CreateBucketOutput, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	if input == nil {
		input = &BucketCreateInput{}
	}

	if input.CreateBucketInput == nil {
		input.CreateBucketInput = &s3.CreateBucketInput{}
	}

	input.Bucket = b.Name

	out, err := b.Client.CreateBucket(context.TODO(), input.CreateBucketInput)

	return out, err
}

type BucketDeleteInput struct {
	*s3.DeleteBucketInput
}

func (b *Bucket) Delete(input *BucketDeleteInput) (*s3.DeleteBucketOutput, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	if input == nil {
		input = &BucketDeleteInput{}
	}

	if input.DeleteBucketInput == nil {
		input.DeleteBucketInput = &s3.DeleteBucketInput{}
	}

	input.Bucket = b.Name

	out, err := b.Client.DeleteBucket(context.TODO(), input.DeleteBucketInput)

	return out, err
}

type BucketUploadObjectInput struct {
	File *[]byte
	Key  *string
	*s3.PutObjectInput
}

func (b *Bucket) UploadObject(input *BucketUploadObjectInput) (*s3.PutObjectOutput, string, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, "", fmt.Errorf("empty 'Name' param")
	}
	if input == nil {
		return nil, "", fmt.Errorf("nil input")
	}
	if *input == (BucketUploadObjectInput{}) {
		return nil, "", fmt.Errorf("empty input")
	}
	if input.File == nil {
		return nil, "", fmt.Errorf("empty 'File' param")
	}
	if input.Key == nil || *input.Key == "" {
		return nil, "", fmt.Errorf("empty 'Key' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, "", err
		}
	}

	if input.PutObjectInput == nil {
		input.PutObjectInput = &s3.PutObjectInput{
			Key: input.Key,
		}
	}

	input.Body = bytes.NewReader(*input.File)
	input.Bucket = b.Name

	out, err := b.Client.PutObject(context.TODO(), input.PutObjectInput)

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *b.Name, *b.Region, *input.Key)

	return out, url, err
}

type BucketGetObjectInput struct {
	Key *string
	*s3.GetObjectInput
}

func (b *Bucket) GetObject(input *BucketGetObjectInput) (*s3.GetObjectOutput, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (BucketGetObjectInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Key == nil || *input.Key == "" {
		return nil, fmt.Errorf("empty 'Key' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	if input.GetObjectInput == nil {
		input.GetObjectInput = &s3.GetObjectInput{
			Key: input.Key,
		}
	}

	input.Bucket = b.Name

	out, err := b.Client.GetObject(context.TODO(), input.GetObjectInput)

	return out, err
}

type BucketDeleteObjectInput struct {
	Key *string
	*s3.DeleteObjectInput
}

func (b *Bucket) DeleteObject(input *BucketDeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (BucketDeleteObjectInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Key == nil || *input.Key == "" {
		return nil, fmt.Errorf("empty 'Key' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	if input.DeleteObjectInput == nil {
		input.DeleteObjectInput = &s3.DeleteObjectInput{
			Key: input.Key,
		}
	}

	input.Bucket = b.Name

	out, err := b.Client.DeleteObject(context.TODO(), input.DeleteObjectInput)

	return out, err
}

type ListObjectsInput struct {
	Prefix     *string
	StartAfter *string
	Limit      *int
	*s3.ListObjectsV2Input
}

func (b *Bucket) ListObjects(input *ListObjectsInput) (*s3.ListObjectsV2Output, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	if input == nil {
		input = &ListObjectsInput{}
	}

	if input.ListObjectsV2Input == nil {
		input.ListObjectsV2Input = &s3.ListObjectsV2Input{}
	}

	input.Bucket = b.Name

	if input.Limit != nil {
		limit := int32(*input.Limit)
		input.MaxKeys = &limit
	}

	out, err := b.Client.ListObjectsV2(context.TODO(), input.ListObjectsV2Input)

	return out, err
}

type PresignGetInput struct {
	Key      *string
	Duration *time.Duration
	*s3.PresignOptions
}

func (b *Bucket) PresignGet(input *PresignGetInput) (*v4.PresignedHTTPRequest, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (PresignGetInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Key == nil || *input.Key == "" {
		return nil, fmt.Errorf("empty 'Key' param")
	}
	if input.Duration == nil {
		return nil, fmt.Errorf("empty 'Duration' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	presignClient := s3.NewPresignClient(b.Client)

	out, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: b.Name,
		Key:    input.Key,
	}, s3.WithPresignExpires(*input.Duration))

	return out, err
}

type PresignPutInput struct {
	Key      *string
	Duration *time.Duration
	*s3.PresignOptions
}

func (b *Bucket) PresignPut(input *PresignPutInput) (*v4.PresignedHTTPRequest, error) {
	if b.Name == nil || *b.Name == "" {
		return nil, fmt.Errorf("empty 'Name' param")
	}
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	if *input == (PresignPutInput{}) {
		return nil, fmt.Errorf("empty input")
	}
	if input.Key == nil || *input.Key == "" {
		return nil, fmt.Errorf("empty 'Key' param")
	}
	if input.Duration == nil {
		return nil, fmt.Errorf("empty 'Duration' param")
	}

	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	presignClient := s3.NewPresignClient(b.Client)

	out, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: b.Name,
		Key:    input.Key,
	}, s3.WithPresignExpires(*input.Duration))

	return out, err
}
