package s3store

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

var awsSession = session.Must(session.NewSession(&aws.Config{
	Credentials:      credentials.NewStaticCredentials("localstack", "localstack", ""),
	Endpoint:         aws.String("http://localhost:4566"),
	Region:           aws.String("eu-west-1"),
	S3ForcePathStyle: aws.Bool(true),
}))

func createSess(purpose string) *store {
	storeInterface := Must(New(awsSession, fmt.Sprintf("dphoto-unit-archive-%s-%s", purpose, time.Now().Format("20060102150405"))))

	adapter := storeInterface.(*store)

	_, err := adapter.s3.CreateBucket(&s3.CreateBucketInput{Bucket: &adapter.bucketName})
	if err != nil {
		panic(err)
	}

	return adapter
}

func TestUpload(t *testing.T) {
	adapter := createSess("upload")

	for _, name := range []string{"/unittest/2021/img-2021-1.jpg", "/unittest/2021/img-2021-1_01.jpg", "/unittest/2021/img-002.jpg"} {
		_, err := adapter.s3.PutObject(&s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader("content of " + name)),
			Bucket: &adapter.bucketName,
			Key:    aws.String(name),
		})
		if err != nil {
			panic(err)
		}
	}

	type args struct {
		key     archive.DestructuredKey
		content io.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"it should upload a file, without special case", args{newKey("unittest/2020/img-2020-1", ".jpg"), strings.NewReader("foo")}, "unittest/2020/img-2020-1.jpg"},
		{"it should find a different key to avoid a clash", args{newKey("unittest/2021/img-002", ".jpg"), strings.NewReader("foo")}, "unittest/2021/img-002_01.jpg"},
		{"it should use a counter to find a different key to not override existing files", args{newKey("unittest/2021/img-2021-1", ".jpg"), strings.NewReader("foo")}, "unittest/2021/img-2021-1_02.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			got, err := adapter.Upload(tt.args.key, tt.args.content)
			if a.NoError(err, tt.name) {
				a.Equal(tt.want, got, tt.name)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	adapter := createSess("copy")

	for _, name := range []string{"unittest/2021/img-2021-1.jpg", "unittest/2021/img-2021-2.jpg"} {
		_, err := adapter.s3.PutObject(&s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader("content of " + name)),
			Bucket: &adapter.bucketName,
			Key:    aws.String(name),
		})
		if err != nil {
			panic(err)
		}
	}

	tests := []struct {
		name    string
		origin  string
		dest    archive.DestructuredKey
		want    string
		wantErr bool
	}{
		{
			name:   "it should copy a file to requested name",
			origin: "/unittest/2021/img-2021-1.jpg",
			dest: archive.DestructuredKey{
				Prefix: "unittest/2021-q1/img-2021-1",
				Suffix: ".jpg",
			},
			want: "unittest/2021-q1/img-2021-1.jpg",
		},
		{
			name:   "it should move a file to an available name",
			origin: "/unittest/2021/img-2021-1.jpg",
			dest: archive.DestructuredKey{
				Prefix: "unittest/2021/img-2021-2",
				Suffix: ".jpg",
			},
			want: "unittest/2021/img-2021-2_01.jpg",
		},
		{
			name:   "it should return an error if the source file doesn't exist",
			origin: "foobar",
			dest: archive.DestructuredKey{
				Prefix: "unittest/2021/img-2021-1",
				Suffix: ".jpg",
			},
			wantErr: true,
		},
		{
			name:   "it should not allow prefix starting wity '/'",
			origin: "foobar",
			dest: archive.DestructuredKey{
				Prefix: "/foobar",
				Suffix: ".jpg",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adapter.Copy(tt.origin, tt.dest)

			if !tt.wantErr && assert.NoError(t, err, tt.name) {
				assert.Equal(t, tt.want, got, tt.name)

				_, err = adapter.s3.GetObject(&s3.GetObjectInput{
					Bucket: &adapter.bucketName,
					Key:    aws.String(got),
				})
				assert.NoError(t, err, tt.name)

			} else if tt.wantErr {
				assert.Error(t, err, tt.name)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	adapter := createSess("delete")

	for _, name := range []string{"unittest/2021/img-2021-1.jpg", "unittest/2021/img-2021-2.jpg", "unittest/2021/img-2021-3.jpg", "unittest/2021/img-2021-4.jpg"} {
		_, err := adapter.s3.PutObject(&s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(strings.NewReader("content of " + name)),
			Bucket: &adapter.bucketName,
			Key:    aws.String(name),
		})
		if err != nil {
			panic(err)
		}
	}

	tests := []struct {
		name string
		ids  []string
	}{
		{"it should delete one file", []string{"unittest/2021/img-2021-1.jpg"}},
		{"it should delete several files", []string{"unittest/2021/img-2021-2.jpg", "unittest/2021/img-2021-3.jpg"}},
		{"it should not fail when file is already gone", []string{"foobar"}},
		{"it should not fail when deleting several files including one that didn't exist", []string{"unittest/2021/img-2021-4.jpg", "foobar"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := adapter.Delete(tt.ids)

			if assert.NoError(t, err, tt.name) {
				for _, key := range tt.ids {
					_, err = adapter.s3.GetObject(&s3.GetObjectInput{
						Bucket: &adapter.bucketName,
						Key:    &key,
					})

					aerr, ok := err.(awserr.Error)
					assert.True(t, ok, tt.name, key, err)
					assert.Equal(t, s3.ErrCodeNoSuchKey, aerr.Code(), tt.name, key)
				}
			}
		})
	}
}

func newKey(prefix, suffix string) archive.DestructuredKey {
	return archive.DestructuredKey{
		Prefix: prefix,
		Suffix: suffix,
	}
}

func Test_store_Get_And_Put(t *testing.T) {
	adapter := createSess("get")

	err := adapter.Put("a/key/that/exists", "hero/avenger", strings.NewReader("I am Ironman"))
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	tests := []struct {
		name            string
		key             string
		wantContent     []byte
		wantLength      int
		wantContentType string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name:            "it should return not found if the object doesn't exist",
			key:             "a/key/that/doesnt/exist",
			wantContent:     nil,
			wantLength:      0,
			wantContentType: "",
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, archive.NotFoundError, err, i)
			},
		},
		{
			name:            "it should return the object with the size and media type",
			key:             "a/key/that/exists",
			wantContent:     []byte("I am Ironman"),
			wantLength:      12,
			wantContentType: "hero/avenger",
			wantErr:         assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReader, gotLength, gotType, err := adapter.Get(tt.key)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.key)) {
				return
			}
			var gotContent []byte
			if gotReader != nil {
				gotContent, err = ioutil.ReadAll(gotReader)
				if !assert.NoError(t, err) {
					return
				}
			}
			assert.Equalf(t, tt.wantContent, gotContent, "Get(%v)", tt.key)
			assert.Equalf(t, tt.wantLength, gotLength, "Get(%v)", tt.key)
			assert.Equalf(t, tt.wantContentType, gotType, "Get(%v)", tt.key)
		})
	}
}
