package kademlia

import "container/list"

const (
	K int = 20
)

type Bucket struct {
	*list.List
}

func (b *Bucket) Update(c Contact) {
	if element, exists := b.Contains(c); exists {
		b.MoveToBack(element)
	} else if b.Len() < K {
		b.PushBack(c)
	} else {
		head := b.Front()
		_, err := b.Head().Ping()

		if err != nil {
			b.Remove(head)
			b.PushBack(c)
		} else {
			b.MoveToBack(head)
		}
	}
}

func (b Bucket) Contains(c Contact) (*list.Element, bool) {
	for e := b.Front(); e != nil; e = e.Next() {
		if e.Value.(Contact).Id.Equals(c.Id) {
			return e, true
		}
	}

	return nil, false
}

func (b Bucket) Slice() []Contact {
	s := []Contact{}

	for e := b.Front(); e != nil; e = e.Next() {
		s = append(s, e.Value.(Contact))
	}

	return s
}

func (b Bucket) Get(index int) (Contact, bool) {
	count := 0

	for e := b.Front(); e != nil; e = e.Next() {
		if count == index {
			return e.Value.(Contact), true
		}

		count++
	}

	return Contact{}, false
}

func (b Bucket) Head() Contact {
	return b.Front().Value.(Contact)
}

func (b Bucket) Tail() Contact {
	return b.Back().Value.(Contact)
}

func NewBucket() Bucket {
	return Bucket{list.New()}
}
