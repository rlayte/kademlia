package kademlia

import (
	"log"
	"net/rpc"
)

type Client struct {
	t *Triplet
	*rpc.Client
}

func HTTPClient(t *Triplet) (c *rpc.Client) {
	log.Println("Connecting to:", t.Address())
	c, err := rpc.DialHTTP("tcp", t.Address())

	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewClient(sender *Triplet, contact *Triplet) Client {
	return Client{sender, HTTPClient(contact)}
}
