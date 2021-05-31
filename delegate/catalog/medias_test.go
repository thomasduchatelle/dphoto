package catalog_test

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRelocateMovedMedias_full(t *testing.T) {
	a := assert.New(t)

	// given
	operator := new(mocks.MoveMediaOperator)
	repository := new(mocks.RepositoryPort)
	catalog.Repository = repository

	const transactionId = "move-transaction-1"

	movedMedias := []*catalog.MovedMedia{
		{Signature: catalog.MediaSignature{}, SourceFolderName: "A", TargetFolderName: "B", Filename: "001"},
		{Signature: catalog.MediaSignature{}, SourceFolderName: "C", TargetFolderName: "B", Filename: "002"},
		{Signature: catalog.MediaSignature{}, SourceFolderName: "A", TargetFolderName: "D", Filename: "003"},
	}

	repository.On("FindReadyMoveTransactions").Return([]*catalog.MoveTransaction{
		{TransactionId: transactionId, Count: 42},
	}, nil)
	repository.On("FindFilesToMove", transactionId, "").Return(movedMedias[:2], "next-page-1", nil)
	repository.On("FindFilesToMove", transactionId, "next-page-1").Return(movedMedias[2:], "", nil)

	repository.On("UpdateMediasLocation", transactionId, movedMedias[:2]).Return(nil)
	repository.On("UpdateMediasLocation", transactionId, movedMedias[2:]).Return(nil)

	operator.On("Continue").Return(true)
	operator.On("UpdateStatus", 0, 42).Return(nil)
	operator.On("UpdateStatus", 2, 42).Return(nil)
	operator.On("UpdateStatus", 3, 42).Return(nil)

	operator.On("Move", catalog.MediaLocation{FolderName: "A", Filename: "001"}, catalog.MediaLocation{FolderName: "B", Filename: "001"}).Return(nil)
	operator.On("Move", catalog.MediaLocation{FolderName: "C", Filename: "002"}, catalog.MediaLocation{FolderName: "B", Filename: "002"}).Return(nil)
	operator.On("Move", catalog.MediaLocation{FolderName: "A", Filename: "003"}, catalog.MediaLocation{FolderName: "D", Filename: "003"}).Return(nil)

	// when
	got, err := catalog.RelocateMovedMedias(operator)

	if a.NoError(err) {
		a.Equal(3, got)
		operator.AssertExpectations(t)
	}
}

func TestRelocateMovedMedias_interrupt(t *testing.T) {
	a := assert.New(t)

	// given
	operator := new(mocks.MoveMediaOperator)
	repository := new(mocks.RepositoryPort)
	catalog.Repository = repository

	const transactionId = "move-transaction-1"

	movedMedias := []*catalog.MovedMedia{
		{Signature: catalog.MediaSignature{}, SourceFolderName: "A", TargetFolderName: "B", Filename: "001"},
		{Signature: catalog.MediaSignature{}, SourceFolderName: "C", TargetFolderName: "B", Filename: "002"},
		{Signature: catalog.MediaSignature{}, SourceFolderName: "A", TargetFolderName: "D", Filename: "003"},
	}

	repository.On("FindReadyMoveTransactions").Return([]*catalog.MoveTransaction{
		{TransactionId: transactionId, Count: 42},
	}, nil)
	repository.On("FindFilesToMove", transactionId, "").Return(movedMedias[:2], "next-page-1", nil)

	repository.On("UpdateMediasLocation", transactionId, movedMedias[:2]).Return(nil)

	operator.On("Continue").Return(true).Once()
	operator.On("Continue").Return(false)
	operator.On("UpdateStatus", 0, 42).Return(nil)
	operator.On("UpdateStatus", 2, 42).Return(nil)

	operator.On("Move", catalog.MediaLocation{FolderName: "A", Filename: "001"}, catalog.MediaLocation{FolderName: "B", Filename: "001"}).Return(nil)
	operator.On("Move", catalog.MediaLocation{FolderName: "C", Filename: "002"}, catalog.MediaLocation{FolderName: "B", Filename: "002"}).Return(nil)

	// when
	got, err := catalog.RelocateMovedMedias(operator)

	if a.NoError(err) {
		a.Equal(2, got)
		operator.AssertExpectations(t)
	}
}

func TestRelocateMovedMedias_empty(t *testing.T) {
	a := assert.New(t)

	// given
	operator := new(mocks.MoveMediaOperator)
	repository := new(mocks.RepositoryPort)
	catalog.Repository = repository

	repository.On("FindReadyMoveTransactions").Return(nil, nil)

	// when
	got, err := catalog.RelocateMovedMedias(operator)

	if a.NoError(err) {
		a.Equal(0, got)
		operator.AssertExpectations(t)
	}
}
