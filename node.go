package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"strings"
)

const (
	A        int = 3
	idLength int = 160
)

type NodeId [idLength / 8]byte

func (id NodeId) String() string {
	return hex.EncodeToString(id[:])
}

func (id NodeId) LessThan(other NodeId) bool {
	for i, b := range id {
		if b != other[i] {
			return b < other[i]
		}
	}

	return false
}

func NewNodeId(address string) (id NodeId) {
	return sha1.Sum([]byte(address))
}

func Xor(current NodeId, target NodeId) (diff NodeId) {
	for i, b := range current {
		diff[i] = b ^ target[i]
	}

	return
}

type Node struct {
	*Triplet
	buckets [idLength]Bucket
}

func (node Node) BucketIndex(id NodeId) int {
	diff := Xor(node.Id, id)

	for byteIndex, b := range diff {
		if b != 0 {
			bits := strconv.FormatInt(int64(b), 2)
			padding := strings.Repeat("0", 8-len(bits))
			bits = padding + bits

			for bitIndex, char := range bits {
				if char != '0' {
					return (byteIndex * 8) + bitIndex
				}
			}
		}
	}

	return idLength - 1
}

func (node *Node) AddToBucket(triplet Triplet) bool {
	bucket := node.ClosestBucket(triplet.Id)
	bucket.Update(triplet)

	return true
}

func (node *Node) Join(seed Triplet) {
	node.AddToBucket(*node.Triplet)
	node.AddToBucket(seed)
}

func (node *Node) Update(t *Triplet) {
	node.AddToBucket(*t)
}

func (node Node) ClosestBucket(target NodeId) Bucket {
	index := node.BucketIndex(target)
	return node.buckets[index]
}

func (node Node) NextBucket(target NodeId) Bucket {
	index := node.BucketIndex(target)

	if index >= idLength-1 {
		return node.buckets[0]
	} else {
		return node.buckets[index+1]
	}
}

func (node Node) PrevBucket(target NodeId) Bucket {
	index := node.BucketIndex(target)

	if index == 0 {
		return node.buckets[idLength-1]
	} else {
		return node.buckets[index-1]
	}
}

func (node *Node) ClosestNodes(target NodeId, quantity int) []Triplet {
	bucket := node.ClosestBucket(target)
	selected := []Triplet{}

	for len(selected) < quantity {
		count := 0

		for bucket.Len() < A {
			if count > idLength/2 {
				break
			}

			bucket.PushBackList(node.NextBucket(target).List)
			bucket.PushBackList(node.PrevBucket(target).List)
			count++
		}

		selected = append(selected, bucket.RandomTriplets(A-len(selected))...)
	}

	return selected
}

func (node *Node) String() string {
	return node.Triplet.String()
}

func NewNode(ip string, port string) (node Node) {
	var buckets [idLength]Bucket

	for i := 0; i < idLength; i++ {
		buckets[i] = NewBucket()
	}

	triplet := Triplet{Ip: ip, Port: port}
	triplet.Id = NewNodeId(triplet.Address())
	node = Node{&triplet, buckets}

	return
}
