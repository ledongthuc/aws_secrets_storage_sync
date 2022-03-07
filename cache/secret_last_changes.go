package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	FileName    string
	LastChanged time.Time
}

type SecretLastChanges struct {
	inner map[string]CacheItem
	l     *sync.RWMutex
}

func NewSecretLastChanges() *SecretLastChanges {
	return &SecretLastChanges{
		inner: map[string]CacheItem{},
		l:     &sync.RWMutex{},
	}
}

func (cache *SecretLastChanges) Get(secretName string) (CacheItem, bool) {
	cache.l.RLock()
	defer cache.l.RUnlock()

	result, existed := cache.inner[secretName]
	return result, existed
}

func (cache *SecretLastChanges) SetTime(secretName string, v time.Time) {
	cache.l.Lock()
	defer cache.l.Unlock()

	item := cache.inner[secretName]
	item.LastChanged = v
	cache.inner[secretName] = item
}

func (cache *SecretLastChanges) Set(secretName string, i CacheItem) {
	cache.l.Lock()
	defer cache.l.Unlock()

	cache.inner[secretName] = i
}
