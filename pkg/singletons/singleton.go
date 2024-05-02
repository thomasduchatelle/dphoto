package singletons

import (
	"fmt"
	"sync"
)

var (
	singletons = make(map[string]any)
	lock       = &sync.Mutex{}
)

func Singleton[S any](newInstance func() (S, error)) (S, error) {
	key := fmt.Sprintf("%T", *new(S))
	if value, exists := singletons[key]; exists {
		return value.(S), nil
	}

	lock.Lock()
	defer lock.Unlock()

	value, err := newInstance()
	singletons[key] = value
	return value, err
}

func MustSingleton[S any](newInstance func() (S, error)) S {
	singleton, err := Singleton(newInstance)
	if err != nil {
		panic(fmt.Sprintf("PANIC - %T couldn't be built: %s", *new(S), err))
	}

	return singleton
}
