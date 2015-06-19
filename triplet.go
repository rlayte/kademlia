package kademlia

import (
	"fmt"
	"log"
)

type Triplet struct {
	Id     NodeId
	Ip     string
	Port   string
	client Client
}

type Contact interface {
	Ping() (PingResponse, error)
}

func (t Triplet) Address() string {
	return fmt.Sprintf("%s:%s", t.Ip, t.Port)
}

func (t Triplet) String() string {
	return fmt.Sprintf("%v --- %s", t.Id, t.Address())
}

func (t Triplet) Ping() (PingResponse, error) {
	request := Request{t.client.t}
	reply := PingResponse{}
	err := t.client.Call("Server.Ping", &request, &reply)

	log.Println("Pong", reply, err)
	return reply, err
}

func NewTriplet(node Node, ip string, port string) Triplet {
	t := Triplet{
		Id:   NewNodeId(ip + ":" + port),
		Ip:   ip,
		Port: port,
	}

	t.client = NewClient(node.Triplet, &t)

	return t
}
