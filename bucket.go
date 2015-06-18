package kademlia

import "container/list"

const (
	K int = 20
)

type Bucket struct {
	*list.List
}

func (b *Bucket) Update(t Contact) {
	if element, exists := b.Contains(t); exists {
		b.MoveToBack(element)
	} else if b.Len() < K {
		b.PushBack(t)
	} else {
		head := b.Front()
		_, err := b.Head().Ping()

		if err != nil {
			b.Remove(head)
			b.PushBack(t)
		} else {
			b.MoveToBack(head)
		}
	}
}

func (b Bucket) Contains(t Contact) (*list.Element, bool) {
	for e := b.Front(); e != nil; e = e.Next() {
		if e.Value == t {
			return e, true
		}
	}

	return nil, false
}

func (b Bucket) Head() Triplet {
	return b.Front().Value.(Triplet)
}

func (b Bucket) Tail() Triplet {
	return b.Back().Value.(Triplet)
}

func NewBucket() Bucket {
	return Bucket{list.New()}
}
