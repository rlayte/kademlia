package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"math"
	"math/rand"
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

func (id NodeId) Equals(other NodeId) bool {
	for i, b := range id {
		if b != other[i] {
			return false
		}
	}

	return true
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
	*Contact
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

func makeShortlist(origin Node, target *Contact) chan Contact {
	shortlist := make(chan Contact, A)
	contacts := origin.ClosestNodes(target.Id, A)
	log.Println("Closest nodes", contacts)

	for _, contact := range contacts {
		shortlist <- contact
	}

	return shortlist
}

func iterateShortlist(shortlist chan Contact, target *Contact, method string) Contact {
	var closestNode Contact
	var closestDistance NodeId
	done := make(chan Contact, 1)
	contacted := map[NodeId]bool{}

	for contact := range shortlist {
		if len(contacted) > K {
			close(shortlist)
			done <- closestNode
		}

		if distance := Xor(target.Id, contact.Id); distance.LessThan(closestDistance) {
			closestDistance = distance
			closestNode = contact
		}

		if _, ok := contacted[contact.Id]; !ok {
			contacted[contact.Id] = true
			response, _ := contact.FindNode(target)

			for _, contact := range response.Contacts {
				shortlist <- contact
			}
		}
	}

	return <-done
}

func nodeLookup(origin Node, target *Contact, method string) Contact {
	shortlist := makeShortlist(origin, target)
	return iterateShortlist(shortlist, target, method)
}

func (node *Node) FindNode(contact *Contact) Contact {
	log.Println("Finding node", contact)
	return nodeLookup(*node, contact, "FindNode")
}

func (node *Node) AddToBucket(contact Contact) bool {
	bucket := node.ClosestBucket(contact.Id)
	bucket.Update(contact)

	return true
}

func (node *Node) Join(seed Contact) {
	node.AddToBucket(seed)
	node.FindNode(node.Contact)
}

func (node *Node) Update(contact *Contact) {
	node.AddToBucket(*contact)
}

func (node Node) ClosestBucket(target NodeId) Bucket {
	index := node.BucketIndex(target)
	return node.buckets[index]
}

func (node Node) NextBucket(index int) Bucket {
	if index >= idLength-1 {
		return node.buckets[index-(idLength-1)]
	} else {
		return node.buckets[index]
	}
}

func (node Node) PrevBucket(index int) Bucket {
	if index < 0 {
		return node.buckets[(idLength-1)+index]
	} else {
		return node.buckets[index]
	}
}

func (node *Node) ClosestNodes(target NodeId, quantity int) (contacts []Contact) {
	index := node.BucketIndex(target)
	bucket := node.ClosestBucket(target)
	items := bucket.Slice()
	chosen := map[int]bool{}
	count := 0

	for len(items) < quantity {
		if count > idLength/2 {
			break
		}

		next := node.NextBucket(index + count)
		prev := node.PrevBucket(index - count)
		items = append(items, next.Slice()...)
		items = append(items, prev.Slice()...)
		count++
	}

	l := int(math.Min(float64(quantity), float64(len(items))))

	for len(contacts) < l {
		index := rand.Intn(len(items))

		if taken := chosen[index]; !taken {
			contacts = append(contacts, items[index])
			chosen[index] = true
		}
	}

	return
}

func (node *Node) String() string {
	return node.Contact.String()
}

func NewNode(ip string, port string) (node Node) {
	var buckets [idLength]Bucket

	for i := 0; i < idLength; i++ {
		buckets[i] = NewBucket()
	}

	node = Node{buckets: buckets}
	node.Contact = &Contact{
		Id:   NewNodeId(ip + ":" + port),
		Ip:   ip,
		Port: port,
	}

	return
}
