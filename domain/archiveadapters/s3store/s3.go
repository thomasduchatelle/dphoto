// Package s3store implements archive.StoreAdapter with AWS S3 backend
package s3store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io"
	"path"
	"strings"
	"time"
)

type StoreAndCache interface {
	archive.StoreAdapter
	archive.CacheAdapter
}

func New(sess *session.Session, bucketName string) (StoreAndCache, error) {
	s3client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &store{
		s3:         s3client,
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
	s3         *s3.S3
	s3Uploader *s3manager.Uploader
	bucketName string
}

func (s *store) Copy(origin string, destination archive.DestructuredKey) (string, error) {
	if strings.HasPrefix(destination.Prefix, "/") {
		return "", errors.Errorf("Prefix must not start with a '/' in key hint %+v", destination)
	}

	destinationKey, err := s.findUniqueFilename(destination)
	if err != nil {
		return "", errors.Wrapf(err, "cannot find a unique destination key")
	}

	_, err = s.s3.CopyObject(&s3.CopyObjectInput{
		Bucket:     &s.bucketName,
		CopySource: aws.String(strings.Trim(path.Join(s.bucketName, origin), "/")),
		Key:        &destinationKey,
	})
	return destinationKey, errors.Wrapf(err, "failed to copy %s -> %s", origin, destinationKey)
}

func (s *store) Delete(keys []string) error {
	for _, key := range keys {
		_, err := s.s3.DeleteObject(&s3.DeleteObjectInput{
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
	object, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})

	if aerr, isAwsError := err.(awserr.Error); isAwsError && aerr.Code() == s3.ErrCodeNoSuchKey {
		return nil, 0, "", archive.NotFoundError
	}

	return object.Body, int(*object.ContentLength), *object.ContentType, errors.Wrapf(err, "couldn't access key %s in %s bucket", key, s.bucketName)
}

func (s *store) Put(key string, mediaType string, content io.Reader) error {
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:        aws.ReadSeekCloser(content),
		Bucket:      &s.bucketName,
		ContentType: &mediaType,
		Key:         &key,
	})
	return errors.Wrapf(err, "failed to PUT %s in bucket %s", key, s.bucketName)
}

func (s *store) SignedURL(key string, duration time.Duration) (string, error) {
	request, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})

	return request.Presign(duration)
}

func (s *store) Upload(keyHint archive.DestructuredKey, content io.Reader) (string, error) {
	if strings.HasPrefix(keyHint.Prefix, "/") {
		return "", errors.Errorf("Prefix must not start with a '/' in key hint %+v", keyHint)
	}

	key, err := s.findUniqueFilename(keyHint)
	if err != nil {
		return "", err
	}

	_, err = s.s3Uploader.Upload(&s3manager.UploadInput{
		Body:   aws.ReadSeekCloser(content),
		Bucket: &s.bucketName,
		Key:    &key,
	})

	return key, errors.Wrapf(err, "upload failed")
}

func (s *store) findUniqueFilename(keyHint archive.DestructuredKey) (string, error) {
	filenames := make(map[string]interface{})

	err := s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: &keyHint.Prefix,
	}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range output.Contents {
			filenames[*obj.Key] = nil
		}
		return true
	})

	candidate := keyHint.Prefix + keyHint.Suffix
	_, clash := filenames[candidate]
	for index := 1; clash; index++ {
		candidate = fmt.Sprintf("%s_%02d%s", keyHint.Prefix, index, keyHint.Suffix)
		_, clash = filenames[candidate]
	}

	return candidate, errors.Wrapf(err, "could not determine a unique name for %s?%s", keyHint.Prefix, keyHint.Suffix)
}

func (s *store) isNotFound(err error) bool {
	aerr, ok := err.(awserr.Error)
	return ok && aerr.Code() == s3.ErrCodeNoSuchKey
}

func (s *store) WalkCacheByPrefix(prefix string, observer func(string)) error {
	return s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: &prefix,
	}, func(output *s3.ListObjectsV2Output, _ bool) bool {
		for _, content := range output.Contents {
			observer(*content.Key)
		}
		return true
	})
}
