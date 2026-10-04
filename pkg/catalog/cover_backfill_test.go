package catalog

//func TestBackfillCovers_BackfillForOwner(t *testing.T) {
//	owner := ownermodel.Owner("ironman")
//	backfillAvengersId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/avengers")}
//	backfillStealthId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName("/stealth")}
//	avengers := &catalog.Album{AlbumId: backfillAvengersId, Name: "Avengers"}
//	stealth := &catalog.Album{AlbumId: backfillStealthId, Name: "Stealth"}
//	avengersCovers := []catalog.Cover{{MediaId: "m-a", Filename: "a.jpg", Origin: catalog.CoverOriginRandom}}
//	stealthCovers := []catalog.Cover{{MediaId: "m-s", Filename: "s.jpg", Origin: catalog.CoverOriginRandom}}
//
//	listFailure := errors.New("list exploded")
//	randomiseFailure := errors.New("randomise exploded")
//	viewFailure := errors.New("view update exploded")
//
//	ownerWithTwoAlbums := func() *findAlbumByOwnerPortFake {
//		return &findAlbumByOwnerPortFake{AlbumsByOwner: map[ownermodel.Owner][]*catalog.Album{
//			owner: {avengers, stealth},
//		}}
//	}
//
//	type fields struct {
//		FindAlbumByOwnerPort      *findAlbumByOwnerPortFake
//		RandomiseCoversPort       *randomiseCoversPortFake
//		BackfillCoversViewUpdater *backfillCoversViewUpdaterFake
//	}
//	type args struct {
//		owner ownermodel.Owner
//	}
//	tests := []struct {
//		name              string
//		fields            fields
//		args              args
//		wantReport        catalog.BackfillReport
//		wantErr           assert.ErrorAssertionFunc
//		expectCalledIds   []catalog.AlbumId
//		expectViewUpdates map[catalog.AlbumId][]catalog.Cover
//	}{
//		{
//			name: "it should randomise every album of the owner and fan the new covers to the view",
//			fields: fields{
//				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
//				RandomiseCoversPort: &randomiseCoversPortFake{
//					Changed: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
//				},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
//			},
//			args:              args{owner: owner},
//			wantReport:        catalog.BackfillReport{Albums: 2},
//			wantErr:           assert.NoError,
//			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
//			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
//		},
//		{
//			name: "it should not update the view when Randomise reports no change",
//			fields: fields{
//				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
//				RandomiseCoversPort: &randomiseCoversPortFake{
//					Changed: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
//				},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
//			},
//			args:              args{owner: owner},
//			wantReport:        catalog.BackfillReport{Albums: 2},
//			wantErr:           assert.NoError,
//			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
//			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
//		},
//		{
//			name: "it should return an empty report when the owner has no album",
//			fields: fields{
//				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{},
//				RandomiseCoversPort:       &randomiseCoversPortFake{},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
//			},
//			args:              args{owner: owner},
//			wantReport:        catalog.BackfillReport{Albums: 0},
//			wantErr:           assert.NoError,
//			expectCalledIds:   nil,
//			expectViewUpdates: nil,
//		},
//		{
//			name: "it should continue after a per-album Randomise failure and skip the view update for that album",
//			fields: fields{
//				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
//				RandomiseCoversPort: &randomiseCoversPortFake{
//					Errors:  map[catalog.AlbumId]error{backfillAvengersId: randomiseFailure},
//					Changed: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
//				},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
//			},
//			args: args{owner: owner},
//			wantReport: catalog.BackfillReport{
//				Albums:   2,
//				Failures: []catalog.BackfillFailure{{AlbumId: backfillAvengersId, Err: randomiseFailure}},
//			},
//			wantErr:           assert.NoError,
//			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
//			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
//		},
//		{
//			name: "it should continue after a per-album view update failure and record it in the report",
//			fields: fields{
//				FindAlbumByOwnerPort: ownerWithTwoAlbums(),
//				RandomiseCoversPort: &randomiseCoversPortFake{
//					Changed: map[catalog.AlbumId][]catalog.Cover{backfillAvengersId: avengersCovers, backfillStealthId: stealthCovers},
//				},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{
//					Errors: map[catalog.AlbumId]error{backfillAvengersId: viewFailure},
//				},
//			},
//			args: args{owner: owner},
//			wantReport: catalog.BackfillReport{
//				Albums:   2,
//				Failures: []catalog.BackfillFailure{{AlbumId: backfillAvengersId, Err: viewFailure}},
//			},
//			wantErr:           assert.NoError,
//			expectCalledIds:   []catalog.AlbumId{backfillAvengersId, backfillStealthId},
//			expectViewUpdates: map[catalog.AlbumId][]catalog.Cover{backfillStealthId: stealthCovers},
//		},
//		{
//			name: "it should return an error when listing the owner's albums fails",
//			fields: fields{
//				FindAlbumByOwnerPort:      &findAlbumByOwnerPortFake{Err: listFailure},
//				RandomiseCoversPort:       &randomiseCoversPortFake{},
//				BackfillCoversViewUpdater: &backfillCoversViewUpdaterFake{},
//			},
//			args:       args{owner: owner},
//			wantReport: catalog.BackfillReport{},
//			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
//				return assert.ErrorIs(t, err, listFailure)
//			},
//			expectCalledIds:   nil,
//			expectViewUpdates: nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			backfill := &catalog.BackfillCovers{
//				FindAlbumByOwnerPort:      tt.fields.FindAlbumByOwnerPort,
//				RandomiseCoversPort:       tt.fields.RandomiseCoversPort,
//				BackfillCoversViewUpdater: tt.fields.BackfillCoversViewUpdater,
//			}
//
//			report, err := backfill.BackfillForOwner(context.Background(), tt.args.owner)
//			if !tt.wantErr(t, err, fmt.Sprintf("BackfillForOwner(%s)", tt.args.owner)) {
//				return
//			}
//			assert.Equal(t, tt.wantReport, report, "backfill report")
//			assert.Equal(t, tt.expectCalledIds, tt.fields.RandomiseCoversPort.Calls, "albums called")
//			assert.Equal(t, tt.expectViewUpdates, tt.fields.BackfillCoversViewUpdater.Updates, "view updates")
//		})
//	}
//}
