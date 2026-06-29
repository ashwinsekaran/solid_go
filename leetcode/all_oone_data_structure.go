package leetcode

// LeetCode 432 - All O'one Data Structure
//
// Pseudo code:
//   Doubly linked list of buckets, each bucket holds a count and a set of keys
//   Map: key -> bucket
//   inc(key): move key to next bucket (count+1), create bucket if needed
//   dec(key): move key to prev bucket (count-1), remove if count==0
//   getMaxKey(): return any key from tail bucket
//   getMinKey(): return any key from head bucket

type bucket struct {
	count    int
	keys     map[string]bool
	prev     *bucket
	next     *bucket
}

type AllOne struct {
	head    *bucket
	tail    *bucket
	keyBucket map[string]*bucket
}

func AllOneConstructor() AllOne {
	head := &bucket{keys: make(map[string]bool)}
	tail := &bucket{keys: make(map[string]bool)}
	head.next = tail
	tail.prev = head
	return AllOne{head: head, tail: tail, keyBucket: make(map[string]*bucket)}
}

func (a *AllOne) insertAfter(b *bucket, count int) *bucket {
	nb := &bucket{count: count, keys: make(map[string]bool), prev: b, next: b.next}
	b.next.prev = nb
	b.next = nb
	return nb
}

func (a *AllOne) remove(b *bucket) {
	b.prev.next = b.next
	b.next.prev = b.prev
}

func (a *AllOne) Inc(key string) {
	if b, ok := a.keyBucket[key]; ok {
		nb := b.next
		if nb == a.tail || nb.count != b.count+1 {
			nb = a.insertAfter(b, b.count+1)
		}
		nb.keys[key] = true
		a.keyBucket[key] = nb
		delete(b.keys, key)
		if len(b.keys) == 0 {
			a.remove(b)
		}
	} else {
		nb := a.head.next
		if nb == a.tail || nb.count != 1 {
			nb = a.insertAfter(a.head, 1)
		}
		nb.keys[key] = true
		a.keyBucket[key] = nb
	}
}

func (a *AllOne) Dec(key string) {
	b := a.keyBucket[key]
	if b.count == 1 {
		delete(a.keyBucket, key)
	} else {
		pb := b.prev
		if pb == a.head || pb.count != b.count-1 {
			pb = a.insertAfter(b.prev, b.count-1)
		}
		pb.keys[key] = true
		a.keyBucket[key] = pb
	}
	delete(b.keys, key)
	if len(b.keys) == 0 {
		a.remove(b)
	}
}

func (a *AllOne) GetMaxKey() string {
	if a.tail.prev == a.head {
		return ""
	}
	for k := range a.tail.prev.keys {
		return k
	}
	return ""
}

func (a *AllOne) GetMinKey() string {
	if a.head.next == a.tail {
		return ""
	}
	for k := range a.head.next.keys {
		return k
	}
	return ""
}
