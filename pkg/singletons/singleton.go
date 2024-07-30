package singletons

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	singletons = make(map[string]any)
)

func Singleton[S any](newInstance func() (S, error)) (S, error) {
	key := fmt.Sprintf("%T", *new(S))
	if key == "<nil>" {
		return *new(S), errors.Errorf("Singleton cannot be used with interfaces, use the explicit SingletonKey method.")
	}

	return SingletonKey(key, newInstance)
}

func SingletonKey[S any](key string, newInstance func() (S, error)) (S, error) {
	if value, exists := singletons[key]; exists {
		return value.(S), nil
	}

	value, err := newInstance()
	singletons[key] = value
	return value, err
}

func MustSingleton[S any](newInstance func() (S, error)) S {
	singleton, err := Singleton(newInstance)
	if err != nil {
		panic(fmt.Sprintf("PANIC - %T couldn't be built or was returned nil: %s", *new(S), err))
	}

	return singleton
}

func MustSingletonKey[S any](key string, newInstance func() (S, error)) S {
	singleton, err := SingletonKey(key, newInstance)
	if err != nil {
		panic(fmt.Sprintf("PANIC - %T couldn't be built or was returned nil: %s", *new(S), err))
	}

	return singleton
}
