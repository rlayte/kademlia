package kademlia

import (
	"log"
	"net/rpc"
)

type Client struct {
	c    *Contact
	http *rpc.Client
}

func HTTPClient(c *Contact) (client *rpc.Client) {
	log.Println("Connecting to:", c.Address())
	client, err := rpc.DialHTTP("tcp", c.Address())

	if err != nil {
		log.Fatal(err)
	}

	return
}

func (c Client) Request(method string, args interface{}, reply interface{}) error {
	if c.http == nil {
		return nil
	} else {
		return c.http.Call(method, args, reply)
	}
}

func NewClient(sender *Contact, contact *Contact) Client {
	return Client{sender, HTTPClient(contact)}
}
