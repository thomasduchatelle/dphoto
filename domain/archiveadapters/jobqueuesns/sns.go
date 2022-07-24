package jobqueuesns

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"path"
)

type MisfiredCacheMessageV1 struct {
	Owner    string `json:"owner"`
	StoreKey string `json:"storeKey"`
	Width    int    `json:"width"`
}

func New(sess *session.Session, topicARN string) archive.JobQueueAdapter {
	return &adapter{
		snsClient: sns.New(sess),
		topicARN:  topicARN,
	}
}

type adapter struct {
	snsClient *sns.SNS
	topicARN  string
}

func (a *adapter) ReportMisfiredCache(owner, storeKey string, width int) error {
	mess, err := json.Marshal(MisfiredCacheMessageV1{
		Owner:    owner,
		StoreKey: storeKey,
		Width:    width,
	})
	if err != nil {
		return errors.Wrapf(err, "serialising SNS message")
	}

	_, err = a.snsClient.Publish(&sns.PublishInput{
		Message: aws.String(string(mess)),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"owner": aStringAttribute(owner),
			"event": aStringAttribute("MisfiredCacheMessageV1"),
		},
		MessageDeduplicationId: aws.String(fmt.Sprintf("misfired-cache-%s", path.Dir(storeKey))),
		MessageGroupId:         &owner,
		TopicArn:               &a.topicARN,
	})
	return err
}

func aStringAttribute(value string) *sns.MessageAttributeValue {
	return &sns.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: &value,
	}
}
