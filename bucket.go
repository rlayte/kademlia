package kademlia

import (
	"container/list"
	"log"
)

const (
	K int = 20
)

type Bucket struct {
	*list.List
}

func (b *Bucket) Update(c Contact) {
	if element, exists := b.Contains(c); exists {
		log.Println("Contact exists", c)
		b.MoveToBack(element)
	} else if b.Len() < K {
		log.Println("Contact doesn't exist. Less than K", c, b.Len())
		b.PushBack(c)
	} else {
		log.Println("Contact doesn't exist. More than K", c, b.Len())
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
