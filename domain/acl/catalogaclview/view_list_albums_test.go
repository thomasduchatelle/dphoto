package catalogaclview_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogaclview"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
)

const (
	pepperUser       = "pepper"
	pepper           = "pepper@stark.com"
	tony             = "tony@stark.com"
	hulk             = "hulk@avenger.com"
	infinityWarAlbum = "InfinityWar"
)

func TestView_ListAlbums(t *testing.T) {
	album1 := &catalog.Album{Owner: pepper, FolderName: "album1"}
	album2 := &catalog.Album{Owner: pepper, FolderName: "album2"}
	tonyAlbum := &catalog.Album{Owner: tony, FolderName: infinityWarAlbum}

	type fields struct {
		UserEmail    string
		CatalogRules func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter)
	}
	tests := []struct {
		name    string
		fields  fields
		filter  catalogaclview.ListAlbumsFilter
		want    []*catalogaclview.AlbumInView
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should mix albums owned, with album shared with user, with indications what shared to whom",
			fields: fields{
				UserEmail: pepperUser,
				CatalogRules: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("Owner").Return(pepper, nil)
					rules.On("SharedWithUserAlbum").Return([]catalog.AlbumId{
						{Owner: tonyAlbum.Owner, FolderName: tonyAlbum.FolderName},
					}, nil)
					rules.On("SharedByUserGrid", pepper).Return(map[string][]string{
						album2.FolderName: {hulk},
						"something/else":  {tony},
					}, nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("FindAllAlbums", pepper).Return([]*catalog.Album{album1, album2}, nil)
					catalogAdapter.On("FindAlbums", []catalog.AlbumId{
						{Owner: tonyAlbum.Owner, FolderName: tonyAlbum.FolderName},
					}).Return([]*catalog.Album{tonyAlbum}, nil)

					return rules, catalogAdapter
				},
			},
			filter: catalogaclview.ListAlbumsFilter{},
			want: []*catalogaclview.AlbumInView{
				{
					Album:    album1,
					SharedTo: nil,
				},
				{
					Album:    album2,
					SharedTo: []string{hulk},
				},
				{
					Album:    tonyAlbum,
					SharedTo: nil,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only take owned albums if filtering out shared ones",
			fields: fields{
				UserEmail: pepperUser,
				CatalogRules: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("Owner").Return(pepper, nil)
					rules.On("SharedByUserGrid", pepper).Return(map[string][]string{
						album2.FolderName: {hulk},
						"something/else":  {tony},
					}, nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("FindAllAlbums", pepper).Return([]*catalog.Album{album1, album2}, nil)

					return rules, catalogAdapter
				},
			},
			filter: catalogaclview.ListAlbumsFilter{OnlyDirectlyOwned: true},
			want: []*catalogaclview.AlbumInView{
				{
					Album:    album1,
					SharedTo: nil,
				},
				{
					Album:    album2,
					SharedTo: []string{hulk},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only have shared albums if no owner",
			fields: fields{
				UserEmail: pepperUser,
				CatalogRules: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("Owner").Return("", nil)
					rules.On("SharedWithUserAlbum").Return([]catalog.AlbumId{
						{Owner: tonyAlbum.Owner, FolderName: tonyAlbum.FolderName},
					}, nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("FindAlbums", []catalog.AlbumId{
						{Owner: tonyAlbum.Owner, FolderName: tonyAlbum.FolderName},
					}).Return([]*catalog.Album{tonyAlbum}, nil)

					return rules, catalogAdapter
				},
			},
			filter: catalogaclview.ListAlbumsFilter{},
			want: []*catalogaclview.AlbumInView{
				{
					Album:    tonyAlbum,
					SharedTo: nil,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should only have owned albums if nothing shared",
			fields: fields{
				UserEmail: pepperUser,
				CatalogRules: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("Owner").Return(pepper, nil)
					rules.On("SharedWithUserAlbum").Return(nil, nil)
					rules.On("SharedByUserGrid", pepper).Return(nil, nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("FindAllAlbums", pepper).Return([]*catalog.Album{album1, album2}, nil)

					return rules, catalogAdapter
				},
			},
			filter: catalogaclview.ListAlbumsFilter{},
			want: []*catalogaclview.AlbumInView{
				{
					Album:    album1,
					SharedTo: nil,
				},
				{
					Album:    album2,
					SharedTo: nil,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "it should return empty if no albums and nothing shared",
			fields: fields{
				UserEmail: pepperUser,
				CatalogRules: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("Owner").Return(pepper, nil)
					rules.On("SharedWithUserAlbum").Return(nil, nil)
					rules.On("SharedByUserGrid", pepper).Return(nil, nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("FindAllAlbums", pepper).Return(nil, nil)

					return rules, catalogAdapter
				},
			},
			filter:  catalogaclview.ListAlbumsFilter{},
			want:    nil,
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, catalogAdapter := tt.fields.CatalogRules(t)
			v := catalogaclview.View{
				UserEmail:      tt.fields.UserEmail,
				CatalogRules:   rules,
				CatalogAdapter: catalogAdapter,
			}
			got, err := v.ListAlbums(tt.filter)

			if tt.wantErr(t, err) && err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
