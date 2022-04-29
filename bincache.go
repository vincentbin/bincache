package main

import (
	"fmt"
	"log"
	"main/cache"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache.Cache
}

var (
	mu       sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("Getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.Cache{CacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (cache.ByteView, error) {
	if key == "" {
		return cache.ByteView{}, fmt.Errorf("key is required")
	}
	ret, ok := g.mainCache.Get(key)
	if !ok {
		return g.load(key)
	}
	log.Println("[BinCache] hit")
	return ret, nil
}

func (g *Group) load(key string) (value cache.ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (cache.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return cache.ByteView{}, err
	}
	value := cache.ByteView{}
	value.B = value.CloneBytes(bytes)
	g.mainCache.Set(key, value)
	return value, nil
}
