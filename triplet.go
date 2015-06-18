package kademlia

import "fmt"

type Triplet struct {
	Id     NodeId
	Ip     string
	Port   string
	client RPC
}

type Contact interface {
	Ping() (Triplet, error)
}

func (t Triplet) Address() string {
	return fmt.Sprintf("%s:%s", t.Ip, t.Port)
}

func (t Triplet) String() string {
	return fmt.Sprintf("%v --- %s", t.Id, t.Address())
}

func (t Triplet) Ping() (Triplet, error) {
	return t.client.Call("Ping", t)
}

func NewTriplet(ip string, port string) Triplet {
	return Triplet{
		Id:     NewNodeId(ip + ":" + port),
		Ip:     ip,
		Port:   port,
		client: Client{},
	}
}
