package kademlia

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func Put(key string, value interface{}) error {
	return nil
}

func Get(key string) (interface{}, error) {
	return true, nil
}

func Join(s Server, ip string, port string) {
	log.Println("Joining network")
	seed := NewContact(s.node, ip, port)
	s.node.Join(seed)
}

func Nodes(s Server) []Contact {
	nodes := []Contact{}

	for _, bucket := range s.node.buckets {
		nodes = append(nodes, bucket.Slice()...)
	}

	return nodes
}

func Start(ip string, port string) Server {
	server := NewServer(ip, port)
	rpc.Register(server)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatal(err)
	}

	go http.Serve(l, nil)

	return server
}
