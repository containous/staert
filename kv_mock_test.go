package staert

import (
	"errors"
	"strings"

	"github.com/docker/libkv/store"
)

// Extremely limited mock store so we can test initialization
type Mock struct {
	Error           bool
	KVPairs         []*store.KVPair
	WatchTreeMethod func() <-chan []*store.KVPair
}

func (s *Mock) Put(key string, value []byte, opts *store.WriteOptions) error {
	s.KVPairs = append(s.KVPairs, &store.KVPair{Key: key, Value: value, LastIndex: 0})
	return nil
}

func (s *Mock) Get(key string, options *store.ReadOptions) (*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("error")
	}
	for _, kvPair := range s.KVPairs {
		if kvPair.Key == key {
			return kvPair, nil
		}
	}
	return nil, nil
}

func (s *Mock) Delete(key string) error {
	return errors.New("delete not supported")
}

// Exists mock
func (s *Mock) Exists(key string, options *store.ReadOptions) (bool, error) {
	return false, errors.New("exists not supported")
}

// Watch mock
func (s *Mock) Watch(key string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan *store.KVPair, error) {
	return nil, errors.New("watch not supported")
}

// WatchTree mock
func (s *Mock) WatchTree(prefix string, stopCh <-chan struct{}, options *store.ReadOptions) (<-chan []*store.KVPair, error) {
	return s.WatchTreeMethod(), nil
}

// NewLock mock
func (s *Mock) NewLock(key string, options *store.LockOptions) (store.Locker, error) {
	return nil, errors.New("NewLock not supported")
}

// List mock
func (s *Mock) List(prefix string, options *store.ReadOptions) ([]*store.KVPair, error) {
	if s.Error {
		return nil, errors.New("error")
	}
	var kv []*store.KVPair
	for _, kvPair := range s.KVPairs {
		if strings.HasPrefix(kvPair.Key, prefix+"/") {
			if secondSlashIndex := strings.IndexRune(kvPair.Key[len(prefix)+1:], '/'); secondSlashIndex == -1 {
				kv = append(kv, kvPair)
			} else {
				dir := &store.KVPair{
					Key: kvPair.Key[:secondSlashIndex+len(prefix)+1],
				}
				kv = append(kv, dir)
			}
		}
	}
	return kv, nil
}

// DeleteTree mock
func (s *Mock) DeleteTree(prefix string) error {
	return errors.New("DeleteTree not supported")
}

// AtomicPut mock
func (s *Mock) AtomicPut(key string, value []byte, previous *store.KVPair, opts *store.WriteOptions) (bool, *store.KVPair, error) {
	return false, nil, errors.New("AtomicPut not supported")
}

// AtomicDelete mock
func (s *Mock) AtomicDelete(key string, previous *store.KVPair) (bool, error) {
	return false, errors.New("AtomicDelete not supported")
}

// Close mock
func (s *Mock) Close() {}
