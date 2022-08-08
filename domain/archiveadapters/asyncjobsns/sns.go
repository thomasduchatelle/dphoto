// Package asyncjobsns use SNS and SQS to queue and process jobs on AWS
package asyncjobsns

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"path"
)

const (
	DefaultImagesPerMessage = 16
)

type WarmUpCacheByFolderMessageV1 struct {
	Owner          string `json:"owner"`
	MissedStoreKey string `json:"missedStoreKey"`
	Width          int    `json:"width"`
}

type ImageToResizeMessageV1 struct {
	Owner    string `json:"owner"`
	MediaId  string `json:"mediaId"`
	StoreKey string `json:"storeKey"`
	Widths   []int  `json:"widths"`
}

func New(sess *session.Session, topicARN string, queueURL string, imagesPerMessage int) archive.AsyncJobAdapter {
	return &adapter{
		snsClient:        sns.New(sess),
		sqsClient:        sqs.New(sess),
		topicARN:         topicARN,
		queueURL:         queueURL,
		imagesPerMessage: imagesPerMessage,
	}
}

type adapter struct {
	snsClient        *sns.SNS
	sqsClient        *sqs.SQS
	topicARN         string
	queueURL         string
	imagesPerMessage int
}

func (a *adapter) WarmUpCacheByFolder(owner, missedStoreKey string, width int) error {
	mess, err := json.Marshal(WarmUpCacheByFolderMessageV1{
		Owner:          owner,
		MissedStoreKey: missedStoreKey,
		Width:          width,
	})
	if err != nil {
		return errors.Wrapf(err, "marshaling [%s, %s, %d]", owner, missedStoreKey, width)
	}

	_, err = a.sqsClient.SendMessage(&sqs.SendMessageInput{
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"ContentType": aSQSStringAttribute("WarmUpCacheByFolderMessageV1"),
		},
		MessageBody:            aws.String(string(mess)),
		MessageDeduplicationId: aws.String(fmt.Sprintf("WarmUpCacheByFolder-%s", path.Dir(missedStoreKey))),
		MessageGroupId:         &owner,
		QueueUrl:               &a.queueURL,
	})
	return errors.Wrapf(err, "sending message to %s : [%s, %s, %d]", a.queueURL, owner, missedStoreKey, width)
}

func (a *adapter) LoadImagesInCache(images ...*archive.ImageToResize) error {
	messageContent := make([]ImageToResizeMessageV1, len(images), len(images))
	for i, img := range images {
		messageContent[i] = ImageToResizeMessageV1{
			Owner:    img.Owner,
			MediaId:  img.MediaId,
			StoreKey: img.StoreKey,
			Widths:   img.Widths,
		}
	}

	for batch := 0; batch < len(messageContent); batch += a.imagesPerMessage {
		end := batch + a.imagesPerMessage
		if end > len(images) {
			end = len(images)
		}
		messageJson, err := json.Marshal(messageContent[batch:end])
		if err != nil {
			return errors.Wrapf(err, "marshaling %d images", len(images))
		}

		_, err = a.snsClient.Publish(&sns.PublishInput{
			Message: aws.String(string(messageJson)),
			MessageAttributes: map[string]*sns.MessageAttributeValue{
				"ContentType": aSNSStringAttribute("[]ImageToResizeMessageV1"),
			},
			TopicArn: &a.topicARN,
		})
		if err != nil {
			return errors.Wrapf(err, "sending to %s SNS %d images", a.topicARN, len(images))
		}
	}

	return nil
}

func aSQSStringAttribute(value string) *sqs.MessageAttributeValue {
	return &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &value,
	}
}
func aSNSStringAttribute(value string) *sns.MessageAttributeValue {
	return &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &value,
	}
}
