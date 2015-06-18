package kademlia

import (
	"log"
	"net/rpc"
)

type RPC interface {
	Call(method string, t Triplet) (Triplet, error)
}

type Client struct {
}

func HTTPClient(t Triplet) (c *rpc.Client) {
	log.Println("Connecting to:", t.Address())
	c, err := rpc.DialHTTP("tcp", t.Address())

	if err != nil {
		log.Fatal(err)
	}

	return
}

func (c Client) Call(method string, t Triplet) (Triplet, error) {
	log.Println("Calling", method)
	httpClient := HTTPClient(t)
	var reply Triplet
	err := httpClient.Call("Server."+method, t, &reply)

	log.Println("Pong", reply)

	return reply, err
}
