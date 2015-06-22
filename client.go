package kademlia

import (
	"log"
	"net/rpc"
)

type Client struct {
	c *Contact
	*rpc.Client
}

func HTTPClient(c *Contact) (client *rpc.Client) {
	log.Println("Connecting to:", c.Address())
	client, err := rpc.DialHTTP("tcp", c.Address())

	if err != nil {
		log.Fatal(err)
	}

	return
}

func NewClient(sender *Contact, contact *Contact) Client {
	return Client{sender, HTTPClient(contact)}
}
