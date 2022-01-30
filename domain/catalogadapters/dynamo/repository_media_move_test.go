package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"sort"
	"testing"
	"time"
)

func setupTest(name string) *Rep {
	suffix := time.Now().Format("20060102150405")

	r := &Rep{
		db: dynamodb.New(session.Must(session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials("localstack", "localstack", ""),
			Endpoint:    aws.String("http://localhost:8000"),
			Region:      aws.String("eu-west-1"),
		}))),
		findMovedMediaBatchSize: 25,
		localDynamodb:           true,
		RootOwner:               "UNITTEST#2",
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
	albums := []catalogmodel.Album{
		{
			Name:       "April 21",
			FolderName: "/media/21-apr",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-01"),
		},
		{
			Name:       "May 21",
			FolderName: "/media/21-may",
			Start:      mustParseDate("2021-05-01"),
			End:        mustParseDate("2021-06-01"),
		},
		{
			Name:       "April Fools Season",
			FolderName: "/media/21-fools",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-05"),
		},
		{
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
	photos := []catalogmodel.CreateMediaRequest{
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img010.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-08"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0010",
				SignatureSize:   1,
			},
		},
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img011.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-09"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0011",
				SignatureSize:   1,
			},
		},
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img012.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-10"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0012",
				SignatureSize:   1,
			},
		},
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img013.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-11"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0013",
				SignatureSize:   1,
			},
		},
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[1].FolderName,
				Filename:   "img014.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-05-01"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0014",
				SignatureSize:   1,
			},
		},
	}
	err := repo.InsertMedias(photos)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// when
	transactionId, count, err := repo.UpdateMedias(
		catalogmodel.NewUpdateFilter().
			WithAlbum(albums[0].FolderName, albums[1].FolderName).
			WithinRange(mustParseDate("2021-04-09"), mustParseDate("2021-04-11")).
			WithinRange(mustParseDate("2021-05-01"), mustParseDate("2021-05-05")),
		albums[2].FolderName,
	)

	// then
	name := "it should move medias within time range from 2 album into a 3rd album"
	a.Equal(3, count, name)
	if a.NoError(err, name) {
		medias, err := repo.FindMedias(albums[0].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err, name) {
			a.Equal([]string{"img010.jpeg", "img013.jpeg"}, extractFilenames("", medias.Content), name)
		}

		medias, err = repo.FindMedias(albums[1].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}

		medias, err = repo.FindMedias(albums[2].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err) {
			a.Equal([]string{"img011.jpeg", "img012.jpeg", "img014.jpeg"}, extractFilenames("", medias.Content), name)
		}

		transactions, err := repo.FindReadyMoveTransactions()
		if a.NoError(err, name) {
			a.Equal([]*catalogmodel.MoveTransaction{
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
	_, _, err = repo.UpdateMedias(catalogmodel.NewUpdateFilter().WithAlbum(albums[2].FolderName), albums[0].FolderName)

	// then
	name = "it should move from an album to another one without range restriction"
	if a.NoError(err, name) {
		medias, err := repo.FindMedias(albums[0].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err, name) {
			a.Equal([]string{"img010.jpeg", "img011.jpeg", "img012.jpeg", "img013.jpeg", "img014.jpeg"}, extractFilenames("", medias.Content), name)
		}

		medias, err = repo.FindMedias(albums[1].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}

		medias, err = repo.FindMedias(albums[2].FolderName, catalogmodel.FindMediaFilter{PageRequest: catalogmodel.PageRequest{}})
		if a.NoError(err) {
			a.Len(medias.Content, 0, name)
		}
	}
}

func TestMoveCycle(t *testing.T) {
	a := assert.New(t)
	repo := setupTest("cycle")

	// given
	albums := []catalogmodel.Album{
		{
			Name:       "April 21",
			FolderName: "/media/21-apr",
			Start:      mustParseDate("2021-04-01"),
			End:        mustParseDate("2021-05-01"),
		},
		{
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
	photos := []catalogmodel.CreateMediaRequest{
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img010.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-08"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0010",
				SignatureSize:   1,
			},
		},
		{
			Location: catalogmodel.MediaLocation{
				FolderName: albums[0].FolderName,
				Filename:   "img011.jpeg",
			},
			Type: "Image",
			Details: catalogmodel.MediaDetails{
				DateTime: mustParseDate("2021-04-09"),
			},
			Signature: catalogmodel.MediaSignature{
				SignatureSha256: "0011",
				SignatureSize:   1,
			},
		},
	}
	err := repo.InsertMedias(photos)
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// and
	transactionId, _, err := repo.UpdateMedias(
		catalogmodel.NewUpdateFilter().WithAlbum(albums[0].FolderName),
		albums[1].FolderName,
	)

	// when
	repo.findMovedMediaBatchSize = 1
	moves, nextPage, err := repo.FindFilesToMove(transactionId, "")
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	err = repo.UpdateMediasLocation(transactionId, moves)

	// then
	if a.NoError(err) {
		a.Len(moves, 1, "it should return 'findMovedMediaBatchSize' moves to perform")

		locations, err := repo.FindMediaLocations(moves[0].Signature)
		name := "it should have updated the current location of the media, and removed the 'to be moved' marker"
		if a.NoError(err, name) {
			a.Equal([]*catalogmodel.MediaLocation{
				{
					FolderName: albums[1].FolderName,
					Filename:   photos[0].Location.Filename,
				},
			}, locations, name)
		}

		locations, err = repo.FindMediaLocations(photos[1].Signature)
		name = "it should have the 2 possible locations of the photo that haven't been physically moved"
		if a.NoError(err, name) {
			a.Equal([]*catalogmodel.MediaLocation{
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
	transactionsBefore, err := repo.FindReadyMoveTransactions()

	transactionId, _, err := repo.startMoveTransaction()
	if !a.NoError(err) {
		a.FailNow(err.Error())
	}

	// when
	err = repo.DeleteEmptyMoveTransaction(transactionId)

	// then
	if a.NoError(err) {
		transactions, err := repo.FindReadyMoveTransactions()
		name := "it should have deleted the transaction now all medias has been moved"
		if a.NoError(err, name) {
			a.Equal(transactionsBefore, transactions, name)
		}
	}
}
