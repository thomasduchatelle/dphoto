package onlinestorage

import (
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
	"time"
)

func TestS3OnlineStorage_UploadFile(t *testing.T) {
	a := assert.New(t)

	s, err := NewS3OnlineStorage("dphoto-unit-"+time.Now().Format("20060102150405"), session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})))
	must(a, err)

	_, err = s.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &s.bucketName})
	must(a, err)

	for _, name := range []string{"/unittest/2021/img-2021-1.jpg", "/unittest/2021/img-2021-1_01.jpg", "/unittest/2021/img-002.jpg"} {
		_, err = s.s3.PutObject(&s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader("content of " + name)),
			Bucket: &s.bucketName,
			Key:    aws.String(name),
		})
		must(a, err)
	}

	type args struct {
		media      backupmodel.ReadableMedia
		folderName string
		filename   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"it should upload a file, without special case", args{newMedia("media-1"), "/2020", "img-2020-1.jpg"}, "img-2020-1.jpg"},
		{"it should find a different key when file already uploaded [memory]", args{newMedia("media-1"), "/2020", "img-2020-1.jpg"}, "img-2020-1_01.jpg"},
		{"it should find a different key again [memory]", args{newMedia("media-1"), "/2020", "img-2020-1.jpg"}, "img-2020-1_02.jpg"},
		{"it should upload a file avoiding a clash with existing file", args{newMedia("media-1"), "/2021", "img-2021-1.jpg"}, "img-2021-1_02.jpg"},
		{"it should upload a file avoiding again a clash", args{newMedia("media-1"), "/2021", "img-2021-1.jpg"}, "img-2021-1_03.jpg"},
	}

	for _, tt := range tests {
		got, err := s.UploadFile("unittest", tt.args.media, tt.args.folderName, tt.args.filename)
		if a.NoError(err, tt.name) {
			a.Equal(tt.want, got, tt.name)
		}
	}
}

func TestMoveFile(t *testing.T) {
	a := assert.New(t)
	
	s, err := NewS3OnlineStorage("dphoto-unit-"+time.Now().Format("20060102150405"), session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String("http://localhost:4566"),
	})))
	must(a, err)

	_, err = s.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &s.bucketName})
	must(a, err)

	name := "it should copy a file and delete the original"

	_, err = s.UploadFile("unittest", newMedia("skywalker"), "jedi", "anakin.jpg")
	if a.NoError(err) {
		_, err = s.MoveFile("unittest", "jedi", "anakin.jpg", "sith")
		a.NoError(err, name)

		_, err = s.s3.GetObject(&s3.GetObjectInput{
			Bucket: &s.bucketName,
			Key:    aws.String("unittest/jedi/anakin.jpg"),
		})
		a.Equal(s3.ErrCodeNoSuchKey, err.(awserr.Error).Code(), name)

		_, err = s.s3.GetObject(&s3.GetObjectInput{
			Bucket: &s.bucketName,
			Key:    aws.String("unittest/sith/anakin.jpg"),
		})
		a.NoError(err, name, name)
	}
}

func must(a *assert.Assertions, err error) {
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}
}

type InMemoryMedia struct {
	content string
}

func (i *InMemoryMedia) SimpleSignature() *backupmodel.SimpleMediaSignature {
	return &backupmodel.SimpleMediaSignature{
		RelativePath: "not-used",
		Size:         uint(len(i.content)),
	}
}

func (i *InMemoryMedia) ReadMedia() (io.Reader, error) {
	return strings.NewReader(i.content), nil
}

func newMedia(content string) backupmodel.ReadableMedia {
	return &InMemoryMedia{
		content: content,
	}
}
