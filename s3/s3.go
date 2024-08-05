package s3

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Bucket struct {
	Name   string `json:"name"`
	Region string `json:"region"`
	Client *s3.Client
}

type NewSessionParams struct {
	Region string
}

func NewSession(params *NewSessionParams) (*s3.Client, error) {
	if params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if *params == (NewSessionParams{}) {
		return nil, fmt.Errorf("empty params")
	}
	if params.Region == "" {
		return nil, fmt.Errorf("missing 'Region' param")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(params.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	svc := s3.NewFromConfig(cfg)

	return svc, nil
}

func (b *Bucket) NewSession() (*s3.Client, error) {
	svc, err := NewSession(&NewSessionParams{
		Region: b.Region,
	})
	if err != nil {
		return nil, err
	}

	b.Client = svc

	return b.Client, nil
}

type ListBucketsParams struct {
	SVC    *s3.Client
	Region string
}

func ListBuckets(params *ListBucketsParams) (*s3.ListBucketsOutput, error) {
	if params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if params.SVC == nil {
		return nil, fmt.Errorf("missing 'SVC' param")
	}
	if params.Region == "" {
		return nil, fmt.Errorf("missing 'Region' param")
	}

	if params.SVC == nil {
		var err error
		params.SVC, err = NewSession(&NewSessionParams{
			Region: params.Region,
		})
		if err != nil {
			return nil, err
		}
	}

	input := &s3.ListBucketsInput{}

	resp, err := params.SVC.ListBuckets(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	return resp, nil
}

type BucketCreateParams struct {
	*s3.CreateBucketInput
}

func (b *Bucket) Create(params *BucketCreateParams) (*s3.CreateBucketOutput, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.CreateBucketInput{}
	if params != nil && params.CreateBucketInput != nil {
		input = params.CreateBucketInput
	}

	input.Bucket = aws.String(b.Name)

	out, err := b.Client.CreateBucket(context.TODO(), input)

	return out, err
}

type BucketDeleteParams struct {
	*s3.DeleteBucketInput
}

func (b *Bucket) Delete(params *BucketDeleteParams) (*s3.DeleteBucketOutput, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.DeleteBucketInput{}
	if params != nil && params.DeleteBucketInput != nil {
		input = params.DeleteBucketInput
	}

	input.Bucket = aws.String(b.Name)

	out, err := b.Client.DeleteBucket(context.TODO(), input)

	return out, err
}

type BucketUploadObjectParams struct {
	File *[]byte
	Key  string
	*s3.PutObjectInput
}

func (b *Bucket) UploadObject(params *BucketUploadObjectParams) (*s3.PutObjectOutput, string, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, "", err
		}
	}

	input := &s3.PutObjectInput{}
	if params != nil && params.PutObjectInput != nil {
		input = params.PutObjectInput
	}

	input.Body = bytes.NewReader(*params.File)
	input.Bucket = &b.Name
	input.Key = &params.Key

	out, err := b.Client.PutObject(context.TODO(), input)

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", b.Name, b.Region, *input.Key)

	return out, url, err
}

type BucketGetObjectParams struct {
	Key string
	*s3.GetObjectInput
}

func (b *Bucket) GetObject(params *BucketGetObjectParams) (*s3.GetObjectOutput, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.GetObjectInput{}
	if params != nil && params.GetObjectInput != nil {
		input = params.GetObjectInput
	}

	input.Bucket = &b.Name
	input.Key = &params.Key

	out, err := b.Client.GetObject(context.TODO(), input)

	return out, err
}

type BucketDeleteObjectParams struct {
	Key string
	*s3.DeleteObjectInput
}

func (b *Bucket) DeleteObject(params *BucketDeleteObjectParams) (*s3.DeleteObjectOutput, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.DeleteObjectInput{}
	if params != nil && params.DeleteObjectInput != nil {
		input = params.DeleteObjectInput
	}

	input.Bucket = aws.String(b.Name)
	input.Key = aws.String(params.Key)

	out, err := b.Client.DeleteObject(context.TODO(), input)

	return out, err
}

type ListObjectsParams struct {
	Prefix string
	Page   string
	Limit  int
	*s3.ListObjectsV2Input
}

func (b *Bucket) ListObjects(params *ListObjectsParams) (*s3.ListObjectsV2Output, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.ListObjectsV2Input{}
	if params != nil && params.ListObjectsV2Input != nil {
		input = params.ListObjectsV2Input
	}

	input.Bucket = aws.String(b.Name)
	if params.Prefix != "" {
		input.Prefix = aws.String(params.Prefix)
	}

	if params.Page != "" {
		input.StartAfter = aws.String(params.Page)
	}

	if params.Limit != 0 {
		limit := int32(params.Limit)

		input.MaxKeys = &limit
	}

	out, err := b.Client.ListObjectsV2(context.TODO(), input)

	return out, err
}

type PresignGetParams struct {
	Key      string
	Duration time.Duration
	*s3.PresignOptions
}

func (b *Bucket) PresignGet(params *PresignGetParams) (*v4.PresignedHTTPRequest, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.PresignOptions{}
	if params != nil && params.PresignOptions != nil {
		input = params.PresignOptions
	}

	input.Expires = params.Duration

	presignClient := s3.NewPresignClient(b.Client)

	out, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(params.Key),
	}, s3.WithPresignExpires(params.Duration))

	return out, err
}

type PresignPutParams struct {
	Key      string
	Duration time.Duration
	*s3.PresignOptions
}

func (b *Bucket) PresignPut(params *PresignPutParams) (*v4.PresignedHTTPRequest, error) {
	if b.Client == nil {
		_, err := b.NewSession()
		if err != nil {
			return nil, err
		}
	}

	input := &s3.PresignOptions{}
	if params != nil && params.PresignOptions != nil {
		input = params.PresignOptions
	}

	input.Expires = params.Duration

	presignClient := s3.NewPresignClient(b.Client)

	out, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.Name),
		Key:    aws.String(params.Key),
	}, s3.WithPresignExpires(params.Duration))

	return out, err
}
