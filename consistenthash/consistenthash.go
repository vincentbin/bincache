package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	nodes    []int
	hashMap  map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(nodes ...string) {
	for _, key := range nodes {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.nodes = append(m.nodes, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.nodes)
}

func (m *Map) Get(key string) string {
	if len(m.nodes) == 0 {
		return ""
	}

	hashVal := int(m.hash([]byte(key)))
	index := sort.Search(len(m.nodes), func(i int) bool {
		return m.nodes[i] >= hashVal
	})
	return m.hashMap[m.nodes[index % len(m.nodes)]]
}