package concurrent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var NoData = fmt.Errorf("no data")

const mod = 10000 * 1

func (t *_Map) hashFunc(i int64) int64 {
	return i % mod
}

type _Map struct {
	lock *sync.Mutex
	root []*node
}

func NewMap() *_Map {
	return &_Map{
		root: make([]*node, mod),
		lock: &sync.Mutex{},
	}
}

func NewMapWithMod(m int64) *_Map {
	if m < 1 {
		panic(fmt.Sprintf("m must be greater then 0"))
	}
	return &_Map{
		lock: &sync.Mutex{},
		root: make([]*node, m),
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func (t *_Map) addNode(n *node, key Comparable, value string) *node {
	if n == nil {
		return createNode(key, value)
	}
	var ret *node
	if key.Less(n.key) == Less {
		ret = createNode(n.key, n.value)
		ret.right = n.right
		ret.left = t.addNode(n.left, key, value)
	} else if key.Less(n.key) == Greater {
		ret = createNode(n.key, n.value)
		ret.left = n.left
		ret.right = t.addNode(n.right, key, value)
	} else {
		ret = createNode(key, value)
		ret.left = n.left
		ret.right = n.right
	}
	updateDepth(ret)
	ret = rotate(ret)
	return ret
}

func (t *_Map) forEach(n *node, walk func(v node)) {
	if n == nil {
		return
	}
	t.forEach(n.left, walk)
	walk(*n)
	t.forEach(n.right, walk)
}

func (t *_Map) Set(key Comparable, value string) {
	t.lock.Lock()
	defer t.lock.Unlock()
	h := t.hashFunc(key.Int64())
	t.root[h] = t.addNode(t.root[h], key, value)
}

func (t *_Map) SetObject(key Comparable, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	t.Set(key, string(data))
	return nil
}

func (t *_Map) get(n *node, key Comparable) (string, bool) {
	if n == nil {
		return "", false
	}
	if key.Less(n.key) == Less {
		return t.get(n.left, key)
	} else if key.Less(n.key) == Greater {
		return t.get(n.right, key)
	} else {
		return n.value, true
	}
}

func (t *_Map) Get(key Comparable) (string, bool) {
	return t.get(t.root[t.hashFunc(key.Int64())], key)
}

func (t *_Map) GetObject(key Comparable, outAddr interface{}) error {
	s, b := t.Get(key)
	if !b || s == "" {
		return NoData
	}
	return json.Unmarshal([]byte(s), outAddr)
}

func (t *_Map) ForEach(walk func(k Comparable, v string, updateTime time.Time)) {
	for i := range t.root {
		t.forEach(t.root[i], func(v node) {
			walk(v.key, v.value, v.updateTime)
		})
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
