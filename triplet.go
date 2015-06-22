package kademlia

import (
	"fmt"
	"log"
)

type Contact struct {
	Id     NodeId
	Ip     string
	Port   string
	client Client
}

type Contactable interface {
	Ping() (PingResponse, error)
}

func (c Contact) Address() string {
	return fmt.Sprintf("%s:%s", c.Ip, c.Port)
}

func (c Contact) String() string {
	return fmt.Sprintf("%v --- %s", c.Id, c.Address())
}

func (c Contact) Ping() (PingResponse, error) {
	request := Request{c.client.c}
	reply := PingResponse{}
	err := c.client.Call("Server.Ping", &request, &reply)

	log.Println("Pong", reply, err)
	return reply, err
}

func NewContact(node Node, ip string, port string) Contact {
	c := Contact{
		Id:   NewNodeId(ip + ":" + port),
		Ip:   ip,
		Port: port,
	}

	c.client = NewClient(node.Contact, &c)

	return c
}
