package kademlia

import (
	"testing"
)

type MockClient struct {
}

func (c MockClient) Call(method string, args interface{}, reply interface{}) error {
	return nil
}

func mockTriplet(id string) Triplet {
	return Triplet{Id: NewNodeId(id), client: MockClient{}}
}

func TestUpdate(t *testing.T) {
	b := NewBucket()
	t1 := mockTriplet("test")

	b.Update(t1)

	if b.Tail().Id != t1.Id {
		t.Error("New nodes should be added to the end of the bucket")
	}

	for i := 1; i < K; i++ {
		b.Update(Triplet{Id: NewNodeId("test" + string(i))})
	}

	t2 := Triplet{Id: NewNodeId("final test")}
	b.Update(t2)

	if b.Tail().Id == t2.Id {
		t.Error("New nodes should not be added if the list is full")
	}

	if b.Len() > K {
		t.Error("Buckets should only contain", K, "items")
	}

	b.Update(t1)

	if b.Tail().Id != t1.Id {
		t.Error("Existing nodes should be moved to the tail when updated")
	}
}

func TestRandomTriplets(t *testing.T) {
	b := NewBucket()

	t1 := mockTriplet("test")
	b.Update(t1)

	random := b.RandomTriplets(3)

	if len(random) != 1 {
		t.Error("Should only return the nodes in the bucket")
	}

	if random[0] != t1 {
		t.Error("Should return the correct triplets")
	}

	for i := 1; i < K; i++ {
		b.Update(mockTriplet("test" + string(i)))
	}

	random = b.RandomTriplets(3)

	if len(random) != 3 {
		t.Error("Should return specified count", 3, len(random))
	}
}
