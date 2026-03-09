package s3

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func New(bucket string) (*S3Storage, error) {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

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

	var err error

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	for i := 0; i < 3; i++ {

		_, err := s.client.PutObject(
			ctx,
			&s3.PutObjectInput{
				Bucket: &s.bucket,
				Key:    &name,
				Body:   r,
			},
		)

		if err != nil {
			return nil
		}
	}

	return err
}
