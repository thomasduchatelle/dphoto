package onlinestorage

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/aws/aws-sdk-go/aws"
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

	for _, name := range []string{"/2021/img-2021-1.jpg", "/2021/img-2021-1_01.jpg", "/2021/img-002.jpg"} {
		_, err = s.s3.PutObject(&s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader("content of " + name)),
			Bucket: &s.bucketName,
			Key:    aws.String(name),
		})
		must(a, err)
	}

	type args struct {
		media      model.ReadableMedia
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
		got, err := s.UploadFile(tt.args.media, tt.args.folderName, tt.args.filename)
		if a.NoError(err, tt.name) {
			a.Equal(tt.want, got, tt.name)
		}
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

func (i *InMemoryMedia) SimpleSignature() *model.SimpleMediaSignature {
	return &model.SimpleMediaSignature{
		RelativePath: "not-used",
		Size:         uint(len(i.content)),
	}
}

func (i *InMemoryMedia) ReadMedia() (io.Reader, error) {
	return strings.NewReader(i.content), nil
}

func newMedia(content string) model.ReadableMedia {
	return &InMemoryMedia{
		content: content,
	}
}
