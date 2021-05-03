package onlinestorage

import (
	"crypto/md5"
	"duchatelle.io/dphoto/dphoto/backup"
	"encoding/hex"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"io"
	"path"
	"strings"
	"sync"
)

type S3OnlineStorage struct {
	s3                *s3.S3
	bucketName        string
	lock              sync.Mutex
	listedFolderNames map[string]interface{}
	prefixes          map[string]interface{}
}

func NewS3OnlineStorage(bucketName string, sess *session.Session) (*S3OnlineStorage, error) {
	s3client := s3.New(sess)

	return &S3OnlineStorage{
		s3:                s3client,
		bucketName:        bucketName,
		lock:              sync.Mutex{},
		listedFolderNames: make(map[string]interface{}),
		prefixes:          make(map[string]interface{}),
	}, nil
}

func Must(storage *S3OnlineStorage, err error) *S3OnlineStorage {
	if err != nil {
		panic(err)
	}

	return storage
}

func (s *S3OnlineStorage) UploadFile(media backup.ReadableMedia, folderName, filename string) (string, error) {
	cleanedFolderName := strings.Trim(folderName, "/")

	prefix, suffix := s.splitKey(cleanedFolderName, filename)
	prefix, err := s.findUniquePrefix(prefix)
	if err != nil {
		return "", err
	}

	key := prefix + suffix

	mediaReader, err := media.ReadMedia()
	if err != nil {
		return "", err
	}

	hash := md5.New()
	teeReader := io.TeeReader(mediaReader, hash)

	output, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(teeReader),
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return "", err
	}

	expected := hex.EncodeToString(hash.Sum(nil))
	if strings.Trim(*output.ETag, "\"") != expected {
		return key, errors.Errorf("Upload failed, invalid MD5 [expected '%s' ; got '%s']", expected, *output.ETag)
	}

	return strings.TrimPrefix(key, cleanedFolderName+"/"), nil
}

func (s *S3OnlineStorage) splitKey(folderName, filename string) (string, string) {
	key := path.Join(folderName, filename)
	suffix := path.Ext(key)

	return strings.TrimPrefix(strings.TrimSuffix(key, suffix), "/"), suffix
}

func (s *S3OnlineStorage) findUniquePrefix(prefix string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	dirPrefix := path.Dir(prefix) + "/"
	if _, listed := s.listedFolderNames[dirPrefix]; !listed {
		err := s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
			Bucket: &s.bucketName,
			Prefix: aws.String(dirPrefix),
		}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
			for _, obj := range output.Contents {
				objPrefix := strings.TrimSuffix(*obj.Key, path.Ext(*obj.Key))
				s.prefixes[objPrefix] = nil
			}
			return true
		})

		if err != nil {
			return "", err
		}

		s.listedFolderNames[dirPrefix] = nil
	}

	candidate := prefix
	_, clash := s.prefixes[candidate]
	index := 1
	for clash {
		candidate = fmt.Sprintf("%s_%02d", prefix, index)
		_, clash = s.prefixes[candidate]
		index++
	}

	s.prefixes[candidate] = nil

	return candidate, nil
}
