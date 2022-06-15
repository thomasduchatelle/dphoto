// Package onlinestorage provides operations to the location where medias are backed-up.
package onlinestorage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

type S3OnlineStorage struct {
	s3                 *s3.S3
	bucketName         string
	lock               sync.Mutex
	filenamesPerFolder map[string]map[string]interface{}
	s3Uploader         *s3manager.Uploader
}

func NewS3OnlineStorage(bucketName string, sess *session.Session) (*S3OnlineStorage, error) {
	s3client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3OnlineStorage{
		s3:                 s3client,
		s3Uploader:         uploader,
		bucketName:         bucketName,
		lock:               sync.Mutex{},
		filenamesPerFolder: make(map[string]map[string]interface{}),
	}, nil
}

func Must(storage *S3OnlineStorage, err error) *S3OnlineStorage {
	if err != nil {
		panic(err)
	}

	return storage
}

func (s *S3OnlineStorage) FetchFile(owner string, folderName, filename string) ([]byte, error) {
	data, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    aws.String(path.Join(owner, strings.Trim(folderName, "/"), filename)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			return nil, backup.MediaNotFoundError
		}
		return nil, err
	}

	defer data.Body.Close()

	return ioutil.ReadAll(data.Body)
}

func (s *S3OnlineStorage) ContentSignedUrl(owner string, folderName string, filename string, duration time.Duration) (string, error) {
	request, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    aws.String(path.Join(owner, strings.Trim(folderName, "/"), filename)),
	})

	return request.Presign(duration)
}

func (s *S3OnlineStorage) UploadFile(owner string, media backup.ReadableMedia, folderName, filename string) (string, error) {
	cleanedFolderName := strings.Trim(folderName, "/")

	key, filename, err := s.findUniqueFilename(owner, cleanedFolderName, filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed getting unique prefix for media %s", media)
	}

	mediaReader, err := media.ReadMedia()
	if err != nil {
		return "", err
	}

	_, err = s.s3Uploader.Upload(&s3manager.UploadInput{
		Body:   aws.ReadSeekCloser(mediaReader),
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return "", errors.Wrapf(err, "uploading %s failed", media)
	}

	return filename, nil
}

func (s *S3OnlineStorage) findUniqueFilename(owner string, folderName string, filename string) (string, string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	dir := path.Join(owner, folderName) + "/"
	filenames, listed := s.filenamesPerFolder[dir]

	if !listed {
		filenames = make(map[string]interface{})
		err := s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
			Bucket: &s.bucketName,
			Prefix: aws.String(dir),
		}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range output.Contents {
				filenames[strings.TrimPrefix(*obj.Key, dir)] = nil
			}
			return true
		})

		if err != nil {
			return "", "", err
		}

		s.filenamesPerFolder[dir] = filenames
	}

	filenameSuffix := path.Ext(filename)
	filenamePrefix := strings.TrimSuffix(filename, filenameSuffix)

	candidate := filenamePrefix + filenameSuffix
	_, clash := filenames[candidate]
	index := 1
	for clash {
		candidate = fmt.Sprintf("%s_%02d%s", filenamePrefix, index, filenameSuffix)
		_, clash = filenames[candidate]
		index++
	}

	filenames[candidate] = nil

	return path.Join(dir, candidate), candidate, nil
}

func (s *S3OnlineStorage) MoveFile(owner string, folderName string, filename string, destFolderName string) (string, error) {
	cleanedFolderName := strings.Trim(destFolderName, "/")

	destKey, filename, err := s.findUniqueFilename(owner, cleanedFolderName, filename)
	if err != nil {
		return "", errors.Wrapf(err, "failed getting unique prefix for media %s/%s", destFolderName, filename)
	}

	origKey := strings.Trim(path.Join(s.bucketName, owner, folderName, filename), "/")

	_, err = s.s3.CopyObject(&s3.CopyObjectInput{
		Bucket:     &s.bucketName,
		CopySource: &origKey,
		Key:        &destKey,
	})
	if err != nil {
		return destKey, errors.Wrapf(err, "failed to copy file %s to %s", origKey, destKey)
	}

	_, err = s.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.bucketName,
		Key:    aws.String(strings.Trim(path.Join(owner, folderName, filename), "/")),
	})
	return filename, errors.Wrapf(err, "failed to remove moved file %s", origKey)
}
