package s3

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *awsS3.Client
	bucket string
}

func New(bucket string) (*S3Storage, error) {

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           "http://localhost:9000",
						SigningRegion: "us-east-1",
					}, nil
				},
			),
		),
	)

	if err != nil {
		return nil, err
	}

	client := awsS3.NewFromConfig(cfg, func(o *awsS3.Options) {
		o.UsePathStyle = true
	})

	return &S3Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Storage) Save(
	ctx context.Context,
	name string,
	r io.Reader,
) error {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {

		_, err := s.client.PutObject(
			ctx,
			&s3.PutObjectInput{
				Bucket: &s.bucket,
				Key:    &name,
				Body:   bytes.NewReader(data),
			},
		)

		if err != nil {
			return nil
		}

		time.Sleep(time.Duration(1<<i) * time.Second)
	}

	return err
}
