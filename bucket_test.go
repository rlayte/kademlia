package kademlia

import (
	"testing"
)

type MockClient struct {
}

func (c MockClient) Call(method string, t Triplet) (Triplet, error) {
	return Triplet{}, nil
}

func mockTriplet(id string) Triplet {
	return Triplet{Id: NewNodeId("test"), client: MockClient{}}
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
