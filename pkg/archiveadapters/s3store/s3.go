// Package s3store implements archive.StoreAdapter with AWS S3 backend
package s3store

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"io"
	"path"
	"strings"
	"time"
)

type StoreAndCache interface {
	archive.StoreAdapter
	archive.CacheAdapter
}

func New(cfg aws.Config, bucketName string, optFns ...func(options *s3.Options)) (StoreAndCache, error) {
	s3client := s3.NewFromConfig(cfg, optFns...)
	presign := s3.NewPresignClient(s3client)
	uploader := manager.NewUploader(s3client)

	return &store{
		client:     s3client,
		presign:    presign,
		s3Uploader: uploader,
		bucketName: bucketName,
	}, nil
}

func Must(storage StoreAndCache, err error) StoreAndCache {
	if err != nil {
		panic(err)
	}

	return storage
}

type store struct {
	client     *s3.Client
	s3Uploader *manager.Uploader
	bucketName string
	presign    *s3.PresignClient
}

func (s *store) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	if strings.HasPrefix(destination.Prefix, "/") {
		return "", errors.Errorf("Prefix must not start with a '/' in key hint %+v", destination)
	}

	destinationKey, err := s.findUniqueFilename(destination)
	if err != nil {
		return "", errors.Wrapf(err, "cannot find a unique destination key")
	}

	_, err = s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     &s.bucketName,
		CopySource: aws.String(strings.Trim(path.Join(s.bucketName, origin), "/")),
		Key:        &destinationKey,
	})
	return destinationKey, errors.Wrapf(err, "failed to copy %s -> %s", origin, destinationKey)
}

func (s *store) Delete(keys []string) error {
	for _, key := range keys {
		_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: &s.bucketName,
			Key:    &key,
		})
		if err != nil && !s.isNotFound(err) {
			return errors.Wrapf(err, "failed to remove moved file %s", key)
		}
	}

	return nil
}

func (s *store) Download(key string) (io.ReadCloser, error) {
	reader, _, _, err := s.Get(key)
	return reader, err
}

func (s *store) Get(key string) (io.ReadCloser, int, string, error) {
	object, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})

	if s.isNotFound(err) {
		return nil, 0, "", archive.NotFoundError
	}
	if err != nil {
		return nil, 0, "", errors.Wrapf(err, "couldn't access key %s in %s bucket", key, s.bucketName)
	}

	return object.Body, int(*object.ContentLength), *object.ContentType, errors.Wrapf(err, "couldn't access key %s in %s bucket", key, s.bucketName)
}

func (s *store) Put(key string, mediaType string, content io.Reader) error {
	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Body:        manager.ReadSeekCloser(content),
		Bucket:      &s.bucketName,
		ContentType: &mediaType,
		Key:         &key,
	})
	return errors.Wrapf(err, "failed to PUT %s in bucket %s", key, s.bucketName)
}

func (s *store) SignedURL(key string, duration time.Duration) (string, error) {
	request, err := s.presign.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	}, func(options *s3.PresignOptions) {
		options.Expires = duration
	})

	return request.URL, err
}

func (s *store) Upload(keyHint archive.DestructuredKey, content io.Reader) (string, error) {
	if strings.HasPrefix(keyHint.Prefix, "/") {
		return "", errors.Errorf("Prefix must not start with a '/' in key hint %+v", keyHint)
	}

	key, err := s.findUniqueFilename(keyHint)
	if err != nil {
		return "", err
	}

	_, err = s.s3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Body:   manager.ReadSeekCloser(content),
		Bucket: &s.bucketName,
		Key:    &key,
	})

	return key, errors.Wrapf(err, "upload failed")
}

func (s *store) findUniqueFilename(keyHint archive.DestructuredKey) (string, error) {
	filenames := make(map[string]interface{})

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: &keyHint.Prefix,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return "", errors.Wrapf(err, "could not determine a unique name for %s?%s: failed listing s3://%s/%s*", keyHint.Prefix, keyHint.Suffix, s.bucketName, keyHint.Prefix)
		}
		for _, obj := range page.Contents {
			filenames[*obj.Key] = nil
		}
	}

	candidate := keyHint.Prefix + keyHint.Suffix
	_, clash := filenames[candidate]
	for index := 1; clash; index++ {
		candidate = fmt.Sprintf("%s_%02d%s", keyHint.Prefix, index, keyHint.Suffix)
		_, clash = filenames[candidate]
	}

	return candidate, nil
}

func (s *store) isNotFound(err error) bool {
	var noSuchKeyErr *types.NoSuchKey
	return errors.As(err, &noSuchKeyErr)
}

func (s *store) WalkCacheByPrefix(prefix string, observer func(string)) error {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: &prefix,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			observer(*obj.Key)
		}
	}

	return nil
}
