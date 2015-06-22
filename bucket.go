package kademlia

import "container/list"

const (
	K int = 20
)

type Bucket struct {
	*list.List
}

func (b *Bucket) Update(c Contactable) {
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

func (b Bucket) Contains(c Contactable) (*list.Element, bool) {
	for e := b.Front(); e != nil; e = e.Next() {
		if e.Value == c {
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

func (b Bucket) RandomContacts(count int) (selected []Contact) {
	indexes := make([]bool, b.Len())
	selectedIndexes := []int{}

	for index, _ := range indexes {
		selectedIndexes = append(selectedIndexes, index)

		if len(selectedIndexes) >= count {
			break
		}
	}

	for _, index := range selectedIndexes {
		node, _ := b.Get(index)
		selected = append(selected, node)
	}

	return
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
