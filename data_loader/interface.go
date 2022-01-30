package data_loader

import (
	"context"
	"time"
)

type Loader func(context.Context) (map[string]Serializable, error)

type Driver interface {
	Store(map[string]Serializable, time.Duration) error
	Fetch(string, Serializable) (Serializable, error)
	All() map[string]Serializable
}

type Interface interface {
	// Get a single value
	Get(string, Serializable) (Serializable, error)
	// All return all data if driver support it
	All() map[string]Serializable
	// Start start the loading process
	Start(context.Context) context.Context
	// Notify is a hack. so we can wait for the first time.
	Notify() <-chan time.Time
}

// RawLoader is a function to handle the loading from any slow source
type RawLoader func(context.Context) (interface{}, error)

type RawDriver interface {
	Store(interface{}, time.Duration) error
	All() interface{}
}

type RawInterface interface {
	// All return all data if driver support it
	All() interface{}
	// Start start the loading process
	Start(context.Context) context.Context
	// Notify is a hack. so we can wait for the first time.
	Notify() <-chan time.Time
}
