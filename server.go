package kademlia

import "log"

type Server struct {
	node Node
}

func nodeLookup(origin Node, target NodeId) Triplet {
	nodes := origin.ClosestNodes(target)
	log.Println("Nodes", nodes)
	return Triplet{}
}

func (s Server) Ping(node *Triplet, reply *Triplet) error {
	log.Println("Ping recieved:", node, s)

	reply.Id = s.node.Id
	reply.Ip = s.node.Ip
	reply.Port = s.node.Port

	s.node.Update(node)

	return nil
}

func Join(s Server, ip string, port string) {
	log.Println("Joining network:", s)
	seed := NewTriplet("0.0.0.0", "3000")
	s.node.Join(seed)
	seed.Ping()
}

func NewServer(ip string, port string) Server {
	return Server{node: NewNode(ip, port)}
}
