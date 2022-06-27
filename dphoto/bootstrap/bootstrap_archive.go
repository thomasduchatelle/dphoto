package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/s3store"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		repositoryAdapter := archivedynamo.Must(archivedynamo.New(cfg.GetAWSSession(), cfg.GetString("catalog.dynamodb.table"), false))
		storeAdapter := s3store.Must(s3store.New(cfg.GetAWSSession(), cfg.GetString("backup.s3.bucket")))
		archive.Init(repositoryAdapter, storeAdapter)
	})
}
