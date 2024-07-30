package pkg

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/filesystemvolume"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestBackupAcceptance(t *testing.T) {
	suite.Run(t, new(BackupTestSuite))
}

type BackupTestSuite struct {
	suite.Suite
	factory    awsfactory.AWSFactory
	owner      ownermodel.Owner
	topMediaId catalog.MediaId
}

func (b *BackupTestSuite) SetupSuite() {
	ctx := context.Background()

	err := initForLocalstack(ctx)
	if !assert.NoError(b.T(), err) {
		return
	}

	owner, err := createRandomUser(ctx)
	if !assert.NoError(b.T(), err) {
		return
	}
	b.owner = owner
}

func (b *BackupTestSuite) Test10_Backup() {
	ctx := context.Background()
	t := b.T()

	volume := filesystemvolume.New("./test_resources/acceptance")
	medias, err := volume.FindMedias()
	if !assert.NoError(t, err, "pre-requisite on source volume") || !assert.Len(t, medias, 3) {
		return
	}

	multiFilesBackup := pkgfactory.NewMultiFilesBackup(ctx)
	report, err := multiFilesBackup(ctx, b.owner, volume,
		backup.WithConcurrentAnalyser(3),
		backup.WithConcurrentCataloguer(3),
		backup.WithConcurrentUploader(3),
		backup.WithBatchSize(1),
	)
	if assert.NoError(t, err) {
		assert.Equal(t, []string{"/2024-Q1"}, report.NewAlbums())
		assert.Equal(t, map[string]*backup.TypeCounter{
			"/2024-Q1": backup.NewTypeCounter(backup.MediaTypeImage, 3, 88185),
		}, report.CountPerAlbum())
	}
}

func (b *BackupTestSuite) Test20_CatalogAndArchive() {
	ctx := context.Background()
	t := b.T()

	albums, err := pkgfactory.AlbumView(ctx).ListAlbums(ctx, usermodel.CurrentUser{
		UserId: usermodel.NewUserId(b.owner.Value()),
		Owner:  &b.owner,
	}, catalogviews.ListAlbumsFilter{})
	if !assert.NoError(t, err) {
		return
	}

	assert.Len(t, albums, 1)
	assert.Equal(t, catalog.FolderName("/2024-Q1"), albums[0].FolderName)
	assert.Equal(t, 3, albums[0].MediaCount)

	albumId := albums[0].AlbumId

	medias, err := pkgfactory.CatalogMediaQueries(ctx).ListMedias(ctx, albumId)
	if !assert.NoError(t, err) {
		return
	}

	assert.Len(t, medias, 3)

	b.topMediaId = medias[0].Id
	url, err := archive.GetMediaOriginalURL(b.owner.Value(), b.topMediaId.Value())
	if !assert.NoError(t, err) {
		return
	}

	resp, err := http.Get(url)
	if !assert.NoError(t, err) {
		return
	}
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
	hash := sha1.New()
	written, err := io.Copy(hash, resp.Body)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, int64(48381), written)
	assert.Equal(t, "f21697f24404920b03ada6bbc9ed61b633584cf0", fmt.Sprintf("%x", hash.Sum(nil)))
}

func (b *BackupTestSuite) Test30_ArchiveCache() {
	t := b.T()

	content, mediaType, err := archive.GetResizedImage(b.owner.Value(), b.topMediaId.Value(), archive.MiniatureCachedWidth, 0)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "image/jpeg", mediaType, "MediaType should be matching the media type.")
	assert.Less(t, len(content), 55*1024, "Miniature should be less than 50KB")
}

func initForLocalstack(ctx context.Context) error {
	// Set the AWS config to use localstack

	factory, err := pkgfactory.StartAWSCloudBuilder(&pkgfactory.StaticAWSAdapterNames{
		DynamoDBNameValue:           "dphoto-local",
		ArchiveMainBucketNameValue:  "dphoto-local",
		ArchiveCacheBucketNameValue: "dphoto-local",
	}).
		OverridesAWSFactory(awsfactory.LocalstackAWSFactory(ctx, awsfactory.LocalstackEndpoint)).
		Build(ctx)
	if err != nil {
		return err
	}

	err = appdynamodb.CreateTableIfNecessary(ctx, pkgfactory.AWSNames.DynamoDBName(), factory.GetDynamoDBClient(), true)
	if err != nil {
		return err
	}

	_, err = factory.GetS3Client().CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(pkgfactory.AWSNames.ArchiveMainBucketName()),
	})
	return err
}

func createRandomUser(ctx context.Context) (ownermodel.Owner, error) {
	timestamp := time.Now().Format("20060102150405")
	owner := ownermodel.Owner(fmt.Sprintf("acceptance+%s@example.com", timestamp))

	repository := pkgfactory.AclRepository(ctx)
	createUser := &aclcore.CreateUser{
		ScopesReader: repository,
		ScopeWriter:  repository,
	}
	err := createUser.CreateUser(owner.Value(), "")

	return owner, err
}
