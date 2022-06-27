// Package s3store implements archive.StoreAdapter with AWS S3 backend
package s3store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io"
	"strings"
)

func New(sess *session.Session, bucketName string) (archive.StoreAdapter, error) {
	s3client := s3.New(sess)
	uploader := s3manager.NewUploader(sess)

	return &store{
		s3:         s3client,
		s3Uploader: uploader,
		bucketName: bucketName,
	}, nil
}

func Must(storage archive.StoreAdapter, err error) archive.StoreAdapter {
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

//func (s *store) FetchFile(owner string, folderName, filename string) ([]byte, error) {
//	data, err := s.s3.GetObject(&s3.GetObjectInput{
//		Bucket: &s.bucketName,
//		Key:    aws.String(path.Join(owner, strings.Trim(folderName, "/"), filename)),
//	})
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
//			return nil, backup.MediaNotFoundError
//		}
//		return nil, err
//	}
//
//	defer data.Body.Close()
//
//	return ioutil.ReadAll(data.Body)
//}
//
//func (s *store) ContentSignedUrl(owner string, folderName string, filename string, duration time.Duration) (string, error) {
//	request, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
//		Bucket: &s.bucketName,
//		Key:    aws.String(path.Join(owner, strings.Trim(folderName, "/"), filename)),
//	})
//
//	return request.Presign(duration)
//}

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

//func (s *store) MoveFile(owner string, folderName string, filename string, destFolderName string) (string, error) {
//	cleanedFolderName := strings.Trim(destFolderName, "/")
//
//	destKey, filename, err := s.findUniqueFilename(owner, cleanedFolderName, filename)
//	if err != nil {
//		return "", errors.Wrapf(err, "failed getting unique prefix for media %s/%s", destFolderName, filename)
//	}
//
//	origKey := strings.Trim(path.Join(s.bucketName, owner, folderName, filename), "/")
//
//	_, err = s.s3.CopyObject(&s3.CopyObjectInput{
//		Bucket:     &s.bucketName,
//		CopySource: &origKey,
//		Key:        &destKey,
//	})
//	if err != nil {
//		return destKey, errors.Wrapf(err, "failed to copy file %s to %s", origKey, destKey)
//	}
//
//	_, err = s.s3.DeleteObject(&s3.DeleteObjectInput{
//		Bucket: &s.bucketName,
//		Key:    aws.String(strings.Trim(path.Join(owner, folderName, filename), "/")),
//	})
//	return filename, errors.Wrapf(err, "failed to remove moved file %s", origKey)
//}
