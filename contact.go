package kademlia

import "fmt"

type Contact struct {
	Id     NodeId
	Ip     string
	Port   string
	client Client
}

type Contactable interface {
	Ping() (PingResponse, error)
	FindNode(contact *Contact) (FindNodeResponse, error)
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
	err := c.client.Request("Server.Ping", &request, &reply)

	return reply, err
}

func (c Contact) FindNode(contact *Contact) (FindNodeResponse, error) {
	request := FindRequest{Sender: c.client.c, Target: contact}
	reply := FindNodeResponse{}

	err := c.client.Request("Server.FindNode", &request, &reply)
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
