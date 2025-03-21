package s3volume

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"io/ioutil"
	"path"
	"testing"
)

func TestShouldFindMediasInS3(t *testing.T) {
	a := assert.New(t)
	ctx := context.Background()

	// given
	factory, err := awsfactory.LocalstackAWSFactory(ctx, awsfactory.LocalstackEndpoint)
	if !assert.NoError(t, err) {
		return
	}
	s3Client := factory.GetS3Client()

	mockBucket := fmt.Sprintf("dphoto-unit-s3source-%s", uuid.Must(uuid.NewUUID()))
	mockFiles := []struct {
		key     string
		content []byte
	}{
		{"image_1.jpg", nil},
		{"my_images/image_2.jpg", []byte("content of image 2")},
		{"my_images/holidays_2022/image_3.JPG", nil},
		{"my_images/video_1.mp4", nil},
		{"my_images_before/image_4.jpg", nil},
	}
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &mockBucket,
	})
	if !assert.NoError(t, err) {
		return
	}
	for _, file := range mockFiles {
		_, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
			Body:   bytes.NewReader(file.content),
			Bucket: &mockBucket,
			Key:    &file.key,
		})
		if !a.NoError(err) {
			a.FailNow("failed to setup test, localstack started?")
		}
	}

	// when
	sourceVolume, err := New(s3Client, fmt.Sprintf("s3://%s/my_images", mockBucket))
	if err != nil {
		a.FailNow(err.Error())
	}

	vol := sourceVolume.(*volume)
	vol.supportedExtensions = map[string]backup.MediaType{
		"jpg": backup.MediaTypeImage,
	}

	medias, err := vol.FindMedias(context.TODO())
	if a.NoError(err, "it should find all medias in s3") {
		found := make([]string, len(medias), len(medias))
		for i, media := range medias {
			found[i] = path.Join(media.MediaPath().Path, media.MediaPath().Filename)
		}

		a.Equal([]string{"holidays_2022/image_3.JPG", "image_2.jpg"}, found, "it should filter out unwanted files to keep only those requested.")

		name := "it should fetch the content of the file"
		contentReader, err := medias[1].ReadMedia()
		if a.NoError(err, name) {
			content, err := ioutil.ReadAll(contentReader)
			if a.NoError(err, name) {
				a.Equal([]byte("content of image 2"), content, name)
			}
		}

		name = "it should get the size of the file"
		a.Equal(18, medias[1].Size(), name)

		name = "it should parse the path into a backup.MediaPath"
		a.Equal(backup.MediaPath{
			ParentFullPath: fmt.Sprintf("s3://%s/my_images", mockBucket),
			Root:           fmt.Sprintf("s3://%s/my_images", mockBucket),
			Path:           "",
			Filename:       "image_2.jpg",
			ParentDir:      "my_images",
		}, medias[1].MediaPath(), name)
		a.Equal(backup.MediaPath{
			ParentFullPath: fmt.Sprintf("s3://%s/my_images/holidays_2022", mockBucket),
			Root:           fmt.Sprintf("s3://%s/my_images", mockBucket),
			Path:           "holidays_2022",
			Filename:       "image_3.JPG",
			ParentDir:      "holidays_2022",
		}, medias[0].MediaPath(), name)

		name = "it should get developer friendly toString name"
		a.Equal(fmt.Sprintf("s3://%s/my_images/holidays_2022/image_3.JPG", mockBucket), medias[0].String())

	}
}
