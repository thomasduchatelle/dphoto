package s3volume

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"io"
	"net/url"
	"path"
	"strings"
	"time"
)

// New creates a new backup.SourceVolume that will find files on an S3 bucket.
func New(client *s3.Client, path string) (backup.SourceVolume, error) {
	return newWithS3Client(client, path)
}

func newWithS3Client(client *s3.Client, path string) (backup.SourceVolume, error) {
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
		s3:                  client,
		bucket:              bucket,
		keyPrefix:           prefix,
		supportedExtensions: backup.SupportedExtensions,
	}, nil
}

type volume struct {
	s3                  *s3.Client
	bucket              string
	keyPrefix           string
	supportedExtensions map[string]backup.MediaType
}

func (s *volume) Children(path backup.MediaPath) (backup.SourceVolume, error) {
	return newWithS3Client(s.s3, path.ParentFullPath)
}

func (s *volume) String() string {
	return fmt.Sprintf("s3://%s%s", s.bucket, s.keyPrefix)
}

func (s *volume) FindMedias() ([]backup.FoundMedia, error) {
	var medias []backup.FoundMedia

	paginator := s3.NewListObjectsV2Paginator(s.s3, &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &s.keyPrefix,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to list medias in %s", s.String())
		}

		for _, obj := range page.Contents {
			ext := strings.TrimPrefix(strings.ToLower(path.Ext(*obj.Key)), ".")
			if _, supported := s.supportedExtensions[ext]; supported {
				medias = append(medias, &s3Media{
					bucket:           s.bucket,
					keyPrefix:        s.keyPrefix,
					s3:               s.s3,
					S3Object:         &obj,
					lastModification: *obj.LastModified,
				})
			}
		}
	}

	return medias, nil
}

type s3Media struct {
	bucket           string
	keyPrefix        string
	s3               *s3.Client
	S3Object         *types.Object
	lastModification time.Time
}

func (s *s3Media) LastModification() time.Time {
	return s.lastModification
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
	object, err := s.s3.GetObject(context.TODO(), &s3.GetObjectInput{
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
