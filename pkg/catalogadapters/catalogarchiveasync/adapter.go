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
	numberOfRelocationPerMessage         = 10
	numberOfMessagesSentToSQSOnEachBatch = 10
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
		capacity := int(math.Ceil(float64(len(transferredIds)) / numberOfRelocationPerMessage))
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
		return processPerBatch(ctx, entries, numberOfMessagesSentToSQSOnEachBatch, a.sendMessageBatch)

	}

	return nil
}

func (a *ArchiveASyncRelocator) sendMessageBatch(ctx context.Context, batch []sqstypes.SendMessageBatchRequestEntry) ([]sqstypes.SendMessageBatchRequestEntry, error) {
	result, err := a.SQSClient.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		Entries:  batch,
		QueueUrl: &a.QueueUrl,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "[ArchiveASyncRelocator] failed to send batched messages to relocate medias to queue %s", a.QueueUrl)
	}

	var failed []sqstypes.SendMessageBatchRequestEntry
	for _, resultEntry := range result.Failed {
		for _, entry := range batch {
			if *entry.Id == *resultEntry.Id {
				failed = append(failed, entry)
				break
			}
		}
	}

	if len(failed) > 0 {
		log.Warnf("[ArchiveASyncRelocator] %d messages to relocate medias failed to be sent to %s, retrying", len(failed), a.QueueUrl)
	}

	return failed, nil
}

func processPerBatch[A any](ctx context.Context, slice []A, batchSize int, process func(context.Context, []A) ([]A, error)) error {
	left := slice
	for len(left) > 0 {
		end := batchSize
		if end > len(left) {
			end = len(left)
		}

		unprocessed, err := process(ctx, left[:end])
		if err != nil {
			return err
		}

		left = append(left[end:], unprocessed...)
	}

	return nil
}
