package catalogaclview_test

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogaclview"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/mocks"
	"testing"
)

func TestView_ListMediasFromAlbum(t *testing.T) {
	nopeError := errors.Errorf("Nope.")
	page := catalog.MediaPage{}

	type fields struct {
		UserEmail string
		mocks     func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter)
	}
	type args struct {
		owner      string
		folderName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *catalog.MediaPage
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "it should return list of media when authorised",
			fields: fields{
				UserEmail: pepperUser,
				mocks: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("CanListMediasFromAlbum", pepper, infinityWarAlbum).Return(nil)

					catalogAdapter := mocks.NewACLViewCatalogAdapter(t)
					catalogAdapter.On("ListMedias", pepper, infinityWarAlbum, mock.Anything).Return(&page, nil)
					return rules, catalogAdapter
				},
			},
			args:    args{owner: pepper, folderName: infinityWarAlbum},
			want:    &page,
			wantErr: assert.NoError,
		},
		{
			name: "it should not list medias if the user is not authorised",
			fields: fields{
				UserEmail: pepperUser,
				mocks: func(t *testing.T) (catalogacl.CatalogRules, catalogaclview.ACLViewCatalogAdapter) {
					rules := mocks.NewCatalogRules(t)
					rules.On("CanListMediasFromAlbum", pepper, infinityWarAlbum).Return(nopeError)

					return rules, mocks.NewACLViewCatalogAdapter(t)
				},
			},
			args: args{owner: pepper, folderName: infinityWarAlbum},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Equal(t, nopeError, err, i)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules, catalogAdapter := tt.fields.mocks(t)
			v := catalogaclview.View{
				UserEmail:      tt.fields.UserEmail,
				CatalogRules:   rules,
				CatalogAdapter: catalogAdapter,
			}

			got, err := v.ListMediasFromAlbum(tt.args.owner, tt.args.folderName)
			if !tt.wantErr(t, err, fmt.Sprintf("ListMediasFromAlbum(%v, %v)", tt.args.owner, tt.args.folderName)) {
				return
			}

			assert.Equalf(t, tt.want, got, "ListMediasFromAlbum(%v, %v)", tt.args.owner, tt.args.folderName)
		})
	}
}
