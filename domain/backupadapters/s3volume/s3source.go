package s3volume

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"io"
	"net/url"
	"path"
	"strings"
)

// New creates a new backup.SourceVolume that will find files on an S3 bucket.
func New(sess *session.Session, path string) (backup.SourceVolume, error) {
	s3Url, err := url.Parse(path)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid s3 path, '%s' must be a valid URL (s3://<bucket-name>[/...])", path)
	}
	if s3Url.Scheme != "s3" {
		return nil, errors.Errorf("invalid s3 path, '%s' scheme must be s3://", path)
	}

	bucket := s3Url.Host
	prefix := strings.Trim(s3Url.Path, "/")
	if prefix != "" {
		prefix += "/"
	}

	return &volume{
		s3:                  s3.New(sess),
		bucket:              bucket,
		keyPrefix:           prefix,
		supportedExtensions: backup.SupportedExtensions,
	}, nil
}

type volume struct {
	s3                  *s3.S3
	bucket              string
	keyPrefix           string
	supportedExtensions map[string]backup.MediaType
}

func (s *volume) String() string {
	return fmt.Sprintf("s3://%s%s", s.bucket, s.keyPrefix)
}

func (s *volume) FindMedias() ([]backup.FoundMedia, error) {
	var medias []backup.FoundMedia

	err := s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &s.keyPrefix,
	}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range output.Contents {
			ext := strings.TrimPrefix(strings.ToLower(path.Ext(*obj.Key)), ".")
			if _, supported := s.supportedExtensions[ext]; supported {
				medias = append(medias, &s3Media{
					bucket:    s.bucket,
					keyPrefix: s.keyPrefix,
					s3:        s.s3,
					S3Object:  obj,
				})
			}
		}
		return true
	})

	return medias, errors.Wrapf(err, "Failed to list medias in %s", s.String())
}

type s3Media struct {
	bucket    string
	keyPrefix string
	s3        *s3.S3
	S3Object  *s3.Object
}

func (s *s3Media) Size() int {
	return int(*s.S3Object.Size)
}

func (s *s3Media) MediaPath() backup.MediaPath {
	bucketSuffix := fmt.Sprintf("s3://%s/", s.bucket)
	pathMiddle := path.Dir(strings.TrimPrefix(*s.S3Object.Key, s.keyPrefix))
	if pathMiddle == "." {
		pathMiddle = ""
	}

	//root := strings.Join([]string{"s3:/", s.bucket, s.keyPrefix}, "/") + "/"
	//fullPath := fmt.Sprintf("s3://%s/%s", s.bucket, *s.S3Object.Key)
	//filename := path.Base(*s.S3Object.Key)
	//
	//relativePath := strings.TrimSuffix(path.Dir(*s.S3Object.Key), "/")
	//
	//parent := path.Base(relativePath)
	//if parent == "" {
	//	parent = s.bucket
	//}

	parentDir := path.Base(path.Dir(*s.S3Object.Key))
	if parentDir == "" {
		parentDir = s.bucket
	}

	return backup.MediaPath{
		ParentFullPath: bucketSuffix + path.Dir(*s.S3Object.Key),
		Root:           strings.TrimSuffix(bucketSuffix+s.keyPrefix, "/"),
		Path:           pathMiddle,
		Filename:       path.Base(*s.S3Object.Key),
		ParentDir:      parentDir,
	}
}

func (s *s3Media) ReadMedia() (io.ReadCloser, error) {
	object, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    s.S3Object.Key,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read S3 object s3://%s/%s", s.bucket, *s.S3Object.Key)
	}

	return object.Body, nil
}

func (s *s3Media) String() string {
	return fmt.Sprintf("s3://%s/%s", s.bucket, *s.S3Object.Key)
}
