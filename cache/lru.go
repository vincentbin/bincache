package cache

import "container/list"

type lru struct {
	maxBytes int64
	usedByte int64
	list     *list.List
	cache    map[string]*list.Element
	CallBack func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

func New(maxBytes int64, callBack func(string, Value)) *lru {
	return &lru{
		maxBytes:  maxBytes,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		CallBack:  callBack,
	}
}

func (c *lru) Get(key string) (value Value, ok bool) {
	element, ok := c.cache[key]
	if ok {
		c.list.MoveToFront(element)
		e := element.Value.(*entry)
		return e.value, true
	}
	return
}

func (c *lru) Set(key string, value Value) {
	element, ok := c.cache[key]
	if ok {
		c.list.MoveToFront(element)
		e := element.Value.(*entry)
		c.usedByte += int64(value.Len() - e.value.Len())
		e.value = value
	} else {
		e := c.list.PushFront(&entry{key, value})
		c.cache[key] = e
		c.usedByte += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedByte {
		c.release()
	}
}

func (c *lru) release() {
	element := c.list.Back()
	if element != nil {
		c.list.Remove(element)
		e := element.Value.(*entry)
		c.usedByte -= int64(len(e.key)) + int64(e.value.Len())
		delete(c.cache, e.key)
		if c.CallBack != nil {
			c.CallBack(e.key, e.value)
		}
	}
}

func (c *lru) Len() int {
	return c.list.Len()
}
