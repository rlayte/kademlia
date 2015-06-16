package kademlia

import (
	"strings"
	"testing"
)

func TestNodeId(t *testing.T) {
	id := NewNodeId()

	if len(id) != 20 {
		t.Error("Id should have 20 bytes", len(id))
	}

	for _, i := range id {
		if i > 255 {
			t.Error("Each byte should be 8 bits", i)
		}
	}
}

func testId(last int) (id NodeId) {
	for i := 0; i < idLength/8-1; i++ {
		id[i] = uint8(1)
	}

	id[idLength/8-1] = uint8(last)

	return
}

func TestXor(t *testing.T) {
	distance := Xor(testId(8), testId(9))
	expected := strings.Repeat("0", 39) + "1"

	if distance[idLength/8-1] != 1 {
		t.Error("n1 and n2 should be 1 apart")
	}

	if distance.String() != expected {
		t.Error("Distance should be", expected, "but was", distance.String())
	}
}

func TestBucketIndex(t *testing.T) {
	n := NewNode("0.0.0.0", 3000)
	id := NewNodeId()

	if index := n.BucketIndex(n.id); index != idLength-1 {
		t.Error("BucketIndex should equal", idLength-1, "not", index)
	}

	id[0] = uint8(0)
	n.id[0] = uint8(255)

	if index := n.BucketIndex(id); index != 0 {
		t.Error("BucketIndex should equal", 0, "not", index)
	}

	id[0] = uint8(5)
	id[1] = uint8(5)
	id[2] = uint8(5)
	id[3] = uint8(9)

	n.id[0] = uint8(5)
	n.id[1] = uint8(5)
	n.id[2] = uint8(5)
	n.id[3] = uint8(7)

	if index := n.BucketIndex(id); index != 28 {
		t.Error("BucketIndex should equal", 28, "not", index)
	}
}

func TestJoin(t *testing.T) {
	n1 := NewNode("0.0.0.0", 3000)
	n2 := NewNode("0.0.0.0", 3001)

	n1.Join(*n2.Triplet)

	bucket := n1.buckets[idLength-1]
	tail := bucket.nodes.Back().Value.(Triplet)

	if tail.id != n1.id {
		t.Error("Node's should add themselves to the final bucket")
	}

	index := n1.BucketIndex(n2.id)
	bucket = n1.buckets[index]
	tail = bucket.nodes.Back().Value.(Triplet)

	if tail.id != n2.id {
		t.Error("The seed node should be added to the correct bucket", n1.buckets[index])
	}
}

func TestUpdate(t *testing.T) {
	n1 := NewNode("0.0.0.0", 3000)
	n2 := NewNode("0.0.0.0", 3001)

	n1.Update(n2.Triplet)

	index := n1.BucketIndex(n2.id)
	bucket := n1.buckets[index]
	tail := bucket.nodes.Back().Value.(Triplet)

	if tail.id != n2.id {
		t.Error("Update should added new nodes to the correct bucket")
	}
}
