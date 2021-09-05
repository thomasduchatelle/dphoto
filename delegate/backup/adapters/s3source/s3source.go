package s3source

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors/analyser"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"io"
	"net/url"
	"path"
	"strings"
	"time"
)

type S3Source struct {
	s3 *s3.S3
}

func NewS3Source(sess *session.Session) *S3Source {
	return &S3Source{
		s3: s3.New(sess),
	}
}

func (s *S3Source) FindMediaRecursively(volume backupmodel.VolumeToBackup, callback func(backupmodel.FoundMedia)) (uint, uint, error) {
	s3Url, err := url.Parse(volume.Path)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "'%s' is not a valid S3 bucket path, it should start with s3://", volume.Path)
	}
	if s3Url.Scheme != "s3" {
		return 0, 0, errors.Errorf("'%s' is not a valid S3 bucket path, it should start with s3://", volume.Path)
	}

	var count, size uint

	bucket := s3Url.Host
	prefix := strings.Trim(s3Url.Path, "/")
	if prefix != "" {
		prefix += "/"
	}

	extensions := analyser.SupportedExtensions

	err = s.s3.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	}, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range output.Contents {
			ext := strings.TrimPrefix(strings.ToLower(path.Ext(*obj.Key)), ".")
			if _, supported := extensions[ext]; supported {
				count++
				size += uint(*obj.Size)
				callback(&s3Media{
					bucket:   bucket,
					s3:       s.s3,
					S3Object: obj,
				})
			}
		}
		return true
	})

	return count, size, errors.Wrapf(err, "Failed to list media from bucket %s (path %s)", bucket, s3Url.Path)
}

type s3Media struct {
	bucket   string
	s3       *s3.S3
	S3Object *s3.Object
}

func (s *s3Media) MediaPath() backupmodel.MediaPath {
	root := fmt.Sprintf("s3://%s/", s.bucket)
	relativePath := strings.TrimSuffix(path.Dir(*s.S3Object.Key), "/")

	parent := path.Base(relativePath)
	if parent == "" {
		parent = s.bucket
	}

	return backupmodel.MediaPath{
		ParentFullPath: root + relativePath,
		Root:           root,
		Path:           relativePath,
		Filename:       path.Base(*s.S3Object.Key),
		ParentDir:      parent,
	}
	// note - requirements for this function is path.Base(...) should have the filename, and 'backup/scan.go' should be able to create the "parent volume"
	//return fmt.Sprintf("s3://%s/%s", s.bucket, strings.TrimPrefix(*s.S3Object.Key, "/"))
}

func (s *s3Media) LastModificationDate() time.Time {
	return *s.S3Object.LastModified
}

func (s *s3Media) SimpleSignature() *backupmodel.SimpleMediaSignature {
	return &backupmodel.SimpleMediaSignature{
		RelativePath: *s.S3Object.Key,
		Size:         uint(*s.S3Object.Size),
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
