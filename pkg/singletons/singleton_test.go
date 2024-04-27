package singletons

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_singleton(t *testing.T) {
	type TestSingleton struct {
		Value string
	}
	resetStorage := func(t *testing.T) {
		clear(singletons)
	}
	const testSingletonKey = "*singletons.TestSingleton"
	newlySingleton := TestSingleton{Value: "newly-generated"}
	existingSingleton := TestSingleton{Value: "existing"}

	type args[S any] struct {
		newInstance func() (S, error)
	}
	type testCase[S any] struct {
		name    string
		before  func(t *testing.T)
		args    args[S]
		want    S
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[*TestSingleton]{
		{
			name:   "it should get a new instance if the store is empty",
			before: resetStorage,
			args: args[*TestSingleton]{newInstance: func() (*TestSingleton, error) {
				return &newlySingleton, nil
			}},
			want:    &newlySingleton,
			wantErr: assert.NoError,
		},
		{
			name: "it should get the existing instance and ignore the newInstance returns",
			before: func(t *testing.T) {
				resetStorage(t)
				singletons[testSingletonKey] = &existingSingleton
			},
			args: args[*TestSingleton]{newInstance: func() (*TestSingleton, error) {
				return &newlySingleton, nil
			}},
			want:    &existingSingleton,
			wantErr: assert.NoError,
		},
		{
			name: "it should get the existing instance and ignore the newInstance error",
			before: func(t *testing.T) {
				resetStorage(t)
				singletons[testSingletonKey] = &existingSingleton
			},
			args: args[*TestSingleton]{newInstance: func() (*TestSingleton, error) {
				return nil, errors.Errorf("TEST an error that should be ignored")
			}},
			want:    &existingSingleton,
			wantErr: assert.NoError,
		},
		{
			name:   "it should forward the error from newSingleton",
			before: resetStorage,
			args: args[*TestSingleton]{newInstance: func() (*TestSingleton, error) {
				return nil, errors.Errorf("TEST an error that should be forwarded")
			}},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(t, err, i) && assert.Equal(t, "TEST an error that should be forwarded", err.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			got, err := Singleton(tt.args.newInstance)

			if !tt.wantErr(t, err, fmt.Sprintf("Singleton()")) {
				return
			}

			if assert.Equalf(t, tt.want, got, "Singleton()") && tt.want != got {
				assert.Failf(t, "pointer is to a different instance", " %s != %s", tt.want, got)
			}
		})
	}
}
