package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"sort"
	"testing"
	"time"
)

const owner = "UNITTEST#3"

func setupTest(name string) *rep {
	suffix := time.Now().Format("20060102150405")

	r := &rep{
		db: dynamodb.New(session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials("localstack", "localstack", ""),
			Endpoint:    aws.String("http://localhost:8000"),
			Region:      aws.String("eu-west-1"),
		}))),
		findMovedMediaBatchSize: 25,
		localDynamodb:           true,
		table:                   fmt.Sprintf("test-medias-move-%s-%s", name, suffix),
	}

	err := r.CreateTableIfNecessary()
	if err != nil {
		panic(err)
	}
	return r
}

func TestUpdateMedias(t *testing.T) {
	a := assert.New(t)
	repo := setupTest("update")

	// given
	albums := []catalog.Album{
		{
			Owner:      owner,
			Name:       "April 21",
			FolderName: "/media/21-apr",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-01"),
		},
		{
			Owner:      owner,
			Name:       "May 21",
			FolderName: "/media/21-may",
			Start:      mustParseDate("2021-05-01"),
			End:        mustParseDate("2021-06-01"),
		},
		{
			Owner:      owner,
			Name:       "April Fools Season",
			FolderName: "/media/21-fools",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-05"),
		},
		{
			Owner:      owner,
			Name:       "Mars 21",
			FolderName: "/media/21-mar",
			Start:      mustParseDate("2021-03-01"),
			End:        mustParseDate("2021-04-01"),
		},
	}
	for _, al := range albums {
		err := repo.InsertAlbum(al)
		if !a.NoError(err) {
			a.FailNow(err.Error())
		}
	}

	// and
	photos := []catalog.CreateMediaRequest{
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img010.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-08"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0010",
				SignatureSize:   1,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img011.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-09"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0011",
				SignatureSize:   1,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img012.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-10"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0012",
				SignatureSize:   1,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img013.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-11"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0013",
				SignatureSize:   1,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: albums[1].FolderName,
				Filename:   "img014.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-05-01"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0014",
				SignatureSize:   1,
			},
		},
	}
	err := repo.InsertMedias(owner, photos)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// when
	transactionId, count, err := repo.UpdateMedias(
		catalog.NewFindMediaRequest(owner).
			WithAlbum(albums[0].FolderName, albums[1].FolderName).
			WithinRange(mustParseDate("2021-04-09"), mustParseDate("2021-04-11")).
			WithinRange(mustParseDate("2021-05-01"), mustParseDate("2021-05-05")),
		albums[2].FolderName,
	)

	// then
	name := "it should move medias within time range from 2 album into a 3rd album"
	a.Equal(3, count, name)
	if a.NoError(err, name) {
		medias, err := repo.FindMedias(owner)
		if a.NoError(err, name) {
			a.Equal([]string{"img010.jpeg", "img013.jpeg"}, extractFilenames("", medias.Content), name)
		}

		medias, err = repo.FindMedias(owner)
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}

		medias, err = repo.FindMedias(owner)
		if a.NoError(err) {
			a.Equal([]string{"img011.jpeg", "img012.jpeg", "img014.jpeg"}, extractFilenames("", medias.Content), name)
		}

		transactions, err := repo.FindReadyMoveTransactions(owner)
		if a.NoError(err, name) {
			a.Equal([]*catalog.MoveTransaction{
				{TransactionId: transactionId, Count: 3},
			}, transactions)
		}

		moveOrders, nextPage, err := repo.FindFilesToMove(transactionId, "")
		if a.NoError(err, name) {
			a.Equal("", nextPage)
			a.Len(moveOrders, 3)

			orders := make([]string, len(moveOrders))
			for i, o := range moveOrders {
				orders[i] = fmt.Sprintf("[%s]%s->%s", o.SourceFilename, o.SourceFolderName, o.TargetFolderName)
			}
			sort.Slice(orders, func(i, j int) bool {
				return orders[i] < orders[j]
			})
			a.Equal([]string{"[img011.jpeg]/media/21-apr->/media/21-fools", "[img012.jpeg]/media/21-apr->/media/21-fools", "[img014.jpeg]/media/21-may->/media/21-fools"}, orders, name)
		}
	}

	// when - without range
	_, _, err = repo.UpdateMedias(catalog.NewFindMediaRequest(owner).WithAlbum(albums[2].FolderName), albums[0].FolderName)

	// then
	name = "it should move from an album to another one without range restriction"
	if a.NoError(err, name) {
		medias, err := repo.FindMedias(owner)
		if a.NoError(err, name) {
			a.Equal([]string{"img010.jpeg", "img011.jpeg", "img012.jpeg", "img013.jpeg", "img014.jpeg"}, extractFilenames("", medias.Content), name)
		}

		medias, err = repo.FindMedias(owner)
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}

		medias, err = repo.FindMedias(owner)
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}
	}
}

func TestMoveCycle(t *testing.T) {
	a := assert.New(t)
	repo := setupTest("cycle")

	// given
	albums := []catalog.Album{
		{
			Owner:      owner,
			Name:       "April 21",
			FolderName: "/media/21-apr",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-01"),
		},
		{
			Owner:      owner,
			Name:       "May 21",
			FolderName: "/media/21-may",
			Start:      mustParseDate("2021-05-01"),
			End:        mustParseDate("2021-06-01"),
		},
	}
	for _, al := range albums {
		err := repo.InsertAlbum(al)
		if !a.NoError(err) {
			a.FailNow(err.Error())
		}
	}

	// and
	photos := []catalog.CreateMediaRequest{
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img010.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-08"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0010",
				SignatureSize:   1,
			},
		},
		{
			Location: catalog.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img011.jpeg",
			},
			Type: "Image",
			Details: catalog.MediaDetails{
				DateTime: mustParseDate("2021-04-09"),
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: "0011",
				SignatureSize:   1,
			},
		},
	}
	err := repo.InsertMedias(owner, photos)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// and
	transactionId, _, err := repo.UpdateMedias(
		catalog.NewFindMediaRequest(owner).WithAlbum(albums[0].FolderName),
		albums[1].FolderName,
	)

	// when
	repo.findMovedMediaBatchSize = 1
	moves, nextPage, err := repo.FindFilesToMove(transactionId, "")
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	err = repo.UpdateMediasLocation(owner, transactionId, moves)

	// then
	if a.NoError(err) {
		a.Len(moves, 1, "it should return 'findMovedMediaBatchSize' moves to perform")

		locations, err := repo.FindMediaLocations(owner, moves[0].Signature)
		name := "it should have updated the current location of the media, and removed the 'to be moved' marker"
		if a.NoError(err, name) {
			a.Equal([]*catalog.MediaLocation{
				{
					FolderName: albums[1].FolderName,
					Filename:   photos[0].Location.Filename,
				},
			}, locations, name)
		}

		locations, err = repo.FindMediaLocations(owner, photos[1].Signature)
		name = "it should have the 2 possible locations of the photo that haven't been physically moved"
		if a.NoError(err, name) {
			a.Equal([]*catalog.MediaLocation{
				{
					FolderName: albums[0].FolderName,
					Filename:   photos[1].Location.Filename,
				},
				{
					FolderName: albums[1].FolderName,
					Filename:   photos[1].Location.Filename,
				},
			}, locations, name)
		}
	}

	// when
	moves, nextPage, err = repo.FindFilesToMove(transactionId, nextPage)

	// then
	name := "it should find the second page even if the first page has been moved and cleanup"
	if a.NoError(err, name) {
		a.Len(moves, 1, name)
		a.Equal(photos[1].Signature, moves[0].Signature, name)
	}
}

func TestRep_DeleteEmptyMoveTransaction(t *testing.T) {
	a := assert.New(t)
	repo := setupTest("cycle")

	// given
	transactionsBefore, err := repo.FindReadyMoveTransactions(owner)

	transactionId, _, err := repo.startMoveTransaction(owner)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// when
	err = repo.DeleteEmptyMoveTransaction(transactionId)

	// then
	if a.NoError(err) {
		transactions, err := repo.FindReadyMoveTransactions(owner)
		name := "it should have deleted the transaction now all medias has been moved"
		if a.NoError(err, name) {
			a.Equal(transactionsBefore, transactions, name)
		}
	}
}