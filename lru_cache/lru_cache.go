package lru_cache

type BiNode[T any] struct {
	value T
	next  *BiNode[T]
	prev  *BiNode[T]
}

type LruCache[K comparable, V any] struct {
	head          *BiNode[V]
	tail          *BiNode[V]
	size          int
	capacity      int
	lookup        map[K]*BiNode[V]
	reverseLookup map[*BiNode[V]]K
}

func NewLruCache[K comparable, V any](capacity int) *LruCache[K, V] {
	return &LruCache[K, V]{capacity: capacity, lookup: make(map[K]*BiNode[V], capacity), head: nil, tail: nil, size: 0, reverseLookup: make(map[*BiNode[V]]K)}
}

func (lru *LruCache[K, V]) Update(key K, value V) {
	node, exists := lru.lookup[key]

	if exists == false {
		node := &BiNode[V]{next: nil, prev: nil, value: value}
		lru.size++
		lru.detach(node)
		lru.prepend(node)
		lru.trimCache()

		lru.lookup[key] = node
		lru.reverseLookup[node] = key
	} else {
		lru.detach(node)
		lru.prepend(node)
		node.value = value
	}
}

func (lru *LruCache[K, V]) Get(key K) (V, bool) {
	node, exists := lru.lookup[key]

	if exists == true {

		lru.detach(node)
		lru.prepend(node)

		return node.value, true
	}

	var r V

	return r, false
}

func (lru *LruCache[K, V]) Clear() {
	for k := range lru.lookup {
		delete(lru.lookup, k)
	}

	for bn := range lru.reverseLookup {
		delete(lru.reverseLookup, bn)
	}

	lru.head = nil
	lru.tail = nil
	lru.size = 0
}

func (lru *LruCache[K, V]) prepend(node *BiNode[V]) {
	if lru.tail == nil {
		lru.head = node
		lru.tail = node
	} else {
		lru.head.prev = node
		node.next = lru.head
		lru.head = node
	}
}

func (lru *LruCache[K, V]) detach(node *BiNode[V]) {
	if node.next != nil {
		node.next.prev = node.prev
	}

	if node.prev != nil {
		node.prev.next = node.next
	}

	if lru.head == node {
		lru.head = node.next
	}

	if lru.tail == node {
		lru.tail = node.prev
	}
}

func (lru *LruCache[K, V]) trimCache() {
	if lru.size <= lru.capacity {
		return
	}

	tail := lru.tail
	lru.detach(tail)
	key := lru.reverseLookup[tail]
	delete(lru.lookup, key)
	delete(lru.reverseLookup, tail)
	lru.size--
}
