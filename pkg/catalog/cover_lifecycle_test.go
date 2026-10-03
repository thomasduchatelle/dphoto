package catalog_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestMediasInsertedCoverObserver_OnMediasInserted(t *testing.T) {
	avengersId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/avengers")}
	stealthId := catalog.AlbumId{Owner: "ironman", FolderName: catalog.NewFolderName("/stealth")}

	type args struct {
		medias map[catalog.AlbumId][]catalog.MediaId
	}
	tests := []struct {
		name                string
		args                args
		expectReconciledIds []catalog.AlbumId
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "it should reconcile the covers of each affected album",
			args: args{medias: map[catalog.AlbumId][]catalog.MediaId{
				avengersId: {"m1", "m2"},
				stealthId:  {"m3"},
			}},
			expectReconciledIds: []catalog.AlbumId{avengersId, stealthId},
			wantErr:             assert.NoError,
		},
		{
			name:                "it should be a no-op when no album received medias",
			args:                args{medias: nil},
			expectReconciledIds: nil,
			wantErr:             assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := &reconcileCoversPortFake{}
			observer := &catalog.MediasInsertedCoverObserver{ReconcileCoversPort: port}

			err := observer.OnMediasInserted(context.Background(), tt.args.medias)
			if !tt.wantErr(t, err, fmt.Sprintf("OnMediasInserted(%v)", tt.args.medias)) {
				return
			}

			var reconciledIds []catalog.AlbumId
			for _, call := range port.Reconciled {
				reconciledIds = append(reconciledIds, call.AlbumId)
			}
			sort.Slice(reconciledIds, func(i, j int) bool {
				return reconciledIds[i].FolderName.String() < reconciledIds[j].FolderName.String()
			})
			assert.Equal(t, tt.expectReconciledIds, reconciledIds, "albums reconciled")

			for _, call := range port.Reconciled {
				assert.Nil(t, call.Added, "Reconcile should be called with no added candidates")
				assert.Nil(t, call.Removed, "Reconcile should be called with no removed medias")
			}
		})
	}
}
