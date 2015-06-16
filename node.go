package kademlia

import (
	"container/list"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
)

const (
	K        int = 20
	A        int = 3
	idLength int = 160
)

type Triplet struct {
	id   NodeId
	ip   string
	port int
}

type Bucket struct {
	nodes *list.List
}

type NodeId [idLength / 8]byte

type Node struct {
	*Triplet
	buckets [idLength]Bucket
}

func (id NodeId) String() string {
	return hex.EncodeToString(id[:])
}

func NewNodeId() (id NodeId) {
	for i := 0; i < idLength/8; i++ {
		id[i] = uint8(rand.Intn(256))
	}

	return
}

func Xor(current NodeId, target NodeId) (diff NodeId) {
	for i, b := range current {
		diff[i] = b ^ target[i]
	}

	return
}

func (node Node) BucketIndex(id NodeId) int {
	diff := Xor(node.id, id)

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
	index := node.BucketIndex(triplet.id)
	bucket := node.buckets[index]
	bucket.nodes.PushBack(triplet)

	return true
}

func (node *Node) Join(seed Triplet) {
	node.AddToBucket(*node.Triplet)
	node.AddToBucket(seed)
}

func (node *Node) Update(t *Triplet) {
	node.AddToBucket(*t)
}

func NewNode(ip string, port int) (node Node) {
	var buckets [idLength]Bucket

	for i := 0; i < idLength; i++ {
		buckets[i] = Bucket{nodes: list.New()}
	}

	triplet := Triplet{id: NewNodeId(), ip: ip, port: port}
	node = Node{&triplet, buckets}

	return
}
