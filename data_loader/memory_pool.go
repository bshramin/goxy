package data_loader

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Serializable interface {
	Encode(io.Writer) error
	Decode(io.Reader) error
}

type memoryPool struct {
	data map[string]Serializable
	lock sync.RWMutex
}

func (m *memoryPool) Store(data map[string]Serializable, t time.Duration) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = data
	return nil
}

func (m *memoryPool) Fetch(key string, data Serializable) (Serializable, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return val, nil
}

func (m *memoryPool) All() map[string]Serializable {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.data
}

// NewMemoryPool return an in memory pool
func NewMemoryPool() Driver {
	return &memoryPool{}
}
