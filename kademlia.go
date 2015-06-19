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
	seed := NewTriplet(s.node, ip, port)
	s.node.Join(seed)
	seed.Ping()
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
