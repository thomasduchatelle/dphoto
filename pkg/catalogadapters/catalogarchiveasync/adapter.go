// Package catalogarchiveasync is queueing messages with SQS
package catalogarchiveasync

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"math"
)

const (
	batchSize = 10
)

type ArchiveASyncRelocator struct {
	SQSClient *sqs.Client
	QueueUrl  string
}

type RelocateMediasDTO struct {
	Owner      string   `json:"owner"`
	FolderName string   `json:"folderName"`
	Ids        []string `json:"ids"`
}

func (a *ArchiveASyncRelocator) OnTransferredMedias(ctx context.Context, transfers catalog.TransferredMedias) error {
	var entries []sqstypes.SendMessageBatchRequestEntry
	for targetAlbumId, transferredIds := range transfers.Transfers {
		capacity := int(math.Ceil(float64(len(transferredIds)) / batchSize))
		ids := make([][]string, capacity)

		for i, id := range transferredIds {
			ids[i%capacity] = append(ids[i%capacity], string(id))
		}

		for _, id := range ids {
			message, err := json.Marshal(RelocateMediasDTO{
				Owner:      targetAlbumId.Owner.String(),
				FolderName: targetAlbumId.FolderName.String(),
				Ids:        id,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to marshal batch message for queue %s", a.QueueUrl)
			}

			entries = append(entries, sqstypes.SendMessageBatchRequestEntry{
				Id:          aws.String(uuid.New().String()),
				MessageBody: aws.String(string(message)),
			})
		}
	}

	for len(entries) > 0 {
		result, err := a.SQSClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			Entries:  entries,
			QueueUrl: &a.QueueUrl,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to send batched messages to relocate medias to queue %s", a.QueueUrl)
		}

		prevEntries := entries
		entries = nil
		for _, resultEntry := range result.Failed {
			for _, entry := range prevEntries {
				if *entry.Id == *resultEntry.Id {
					entries = append(entries, entry)
					break
				}
			}
		}

		if len(entries) > 0 {
			log.Warn("%d messages failed to be sent to %s, retrying", len(entries), a.QueueUrl)
		}
	}

	return nil
}
