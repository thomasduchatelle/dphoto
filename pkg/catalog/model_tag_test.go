package catalog_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func TestNewTag(t *testing.T) {
	isEmptyTagError := func(t assert.TestingT, err error, i ...interface{}) bool {
		return assert.ErrorIs(t, err, catalog.EmptyTagError)
	}

	type args struct {
		raw string
	}
	tests := []struct {
		name          string
		args          args
		wantNamespace string
		wantValue     string
		wantString    string
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name:          "it should parse a namespaced tag into its namespace and value",
			args:          args{raw: "people:luke"},
			wantNamespace: "people",
			wantValue:     "luke",
			wantString:    "people:luke",
			wantErr:       assert.NoError,
		},
		{
			name:          "it should round-trip the reserved favorite tag as a bare tag",
			args:          args{raw: "favorite"},
			wantNamespace: "",
			wantValue:     "favorite",
			wantString:    "favorite",
			wantErr:       assert.NoError,
		},
		{
			name:          "it should lowercase an uppercased tag so People:Luke and people:luke are the same",
			args:          args{raw: "People:Luke"},
			wantNamespace: "people",
			wantValue:     "luke",
			wantString:    "people:luke",
			wantErr:       assert.NoError,
		},
		{
			name:          "it should trim surrounding spaces and remove inner spaces",
			args:          args{raw: "  place : New York  "},
			wantNamespace: "place",
			wantValue:     "newyork",
			wantString:    "place:newyork",
			wantErr:       assert.NoError,
		},
		{
			name:          "it should split on the first colon only",
			args:          args{raw: "activity:ski:downhill"},
			wantNamespace: "activity",
			wantValue:     "ski:downhill",
			wantString:    "activity:ski:downhill",
			wantErr:       assert.NoError,
		},
		{
			name:          "it should strip URL-unsafe characters to keep the tag usable in a path",
			args:          args{raw: "people:luke/skywalker?"},
			wantNamespace: "people",
			wantValue:     "lukeskywalker",
			wantString:    "people:lukeskywalker",
			wantErr:       assert.NoError,
		},
		{
			name:    "it should reject a blank tag",
			args:    args{raw: "   "},
			wantErr: isEmptyTagError,
		},
		{
			name:    "it should reject a namespaced tag with an empty value",
			args:    args{raw: "people:"},
			wantErr: isEmptyTagError,
		},
		{
			name:    "it should reject a namespaced tag with an empty namespace",
			args:    args{raw: ":luke"},
			wantErr: isEmptyTagError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := catalog.NewTag(tt.args.raw)
			if !tt.wantErr(t, err) || err != nil {
				return
			}
			assert.Equal(t, tt.wantNamespace, got.Namespace(), "Namespace()")
			assert.Equal(t, tt.wantValue, got.Value(), "Value()")
			assert.Equal(t, tt.wantString, got.String(), "String()")
			assert.Equal(t, mustTag(tt.wantString), got, "the whole tag")
		})
	}
}
