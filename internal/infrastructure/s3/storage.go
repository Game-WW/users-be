package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/minio-go/v7"
)

var ErrObjectNotFound = errors.New("s3 object not found")

type ObjectMeta struct {
	ContentType string
	Size        int64
}

type Storage struct {
	client *minio.Client
	bucket string
}

func NewStorage(client *minio.Client, bucket string) *Storage {
	return &Storage{
		client: client,
		bucket: bucket,
	}
}

func (s *Storage) Put(ctx context.Context, key string, data []byte, contentType string) error {
	reader := bytes.NewReader(data)
	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := s.client.PutObject(ctx, s.bucket, key, reader, int64(len(data)), opts)
	if err != nil {
		return fmt.Errorf("put object %q: %w", key, err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, key string) (io.ReadCloser, ObjectMeta, error) {
	opts := minio.GetObjectOptions{}
	obj, err := s.client.GetObject(ctx, s.bucket, key, opts)
	if err != nil {
		return nil, ObjectMeta{}, fmt.Errorf("get object %q: %w", key, err)
	}

	info, err := obj.Stat()
	if err != nil {
		_ = obj.Close()

		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return nil, ObjectMeta{}, ErrObjectNotFound
		}

		return nil, ObjectMeta{}, fmt.Errorf("stat object %q: %w", key, err)
	}

	meta := ObjectMeta{
		ContentType: info.ContentType,
		Size:        info.Size,
	}

	return obj, meta, nil
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.client.RemoveObject(ctx, s.bucket, key, opts)
	if err != nil {
		return fmt.Errorf("remove object %q: %w", key, err)
	}

	return nil
}

func (s *Storage) ListKeys(ctx context.Context, prefix string) ([]string, error) {
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	var keys []string
	for obj := range s.client.ListObjects(ctx, s.bucket, opts) {
		if obj.Err != nil {
			return nil, fmt.Errorf("list objects with prefix %q: %w", prefix, obj.Err)
		}

		if strings.HasSuffix(obj.Key, "/") {
			continue
		}

		keys = append(keys, obj.Key)
	}

	return keys, nil
}
