package kademlia

import "log"

type Server struct {
	node Node
}

type Request struct {
	Sender *Triplet
}

type FindRequest struct {
	*Request
	Target *Triplet
}

type PingResponse struct {
	*Triplet
}

type FindNodeResponse struct {
	Triplets []Triplet
}

func addToShortlist(shortlist chan Triplet, items []Triplet) {
	for _, triplet := range items {
		shortlist <- triplet
	}
}

func requestClosest(t Triplet, method string) []Triplet {
	reply := FindNodeResponse{}
	t.client.Call(method, &t, reply)

	return reply.Triplets
}

func iterateShortlist(shortlist chan Triplet, target NodeId, method string) Triplet {
	var closestNode Triplet
	var closestDistance NodeId
	done := make(chan Triplet, 1)
	contacted := map[NodeId]bool{}

	for triplet := range shortlist {
		if len(contacted) > K {
			close(shortlist)
			done <- closestNode
		}

		if distance := Xor(target, triplet.Id); distance.LessThan(closestDistance) {
			closestDistance = distance
			closestNode = triplet
		}

		if _, ok := contacted[triplet.Id]; !ok {
			contacted[triplet.Id] = true
			addToShortlist(shortlist, requestClosest(triplet, method))
		}
	}

	return <-done
}

func nodeLookup(origin Node, target NodeId, method string) Triplet {
	shortlist := make(chan Triplet, A)
	addToShortlist(shortlist, origin.ClosestNodes(target, A))
	return iterateShortlist(shortlist, target, method)
}

func (s Server) Ping(request *Request, reply *PingResponse) error {
	log.Println("Ping recieved:", request.Sender)

	reply.Triplet = &Triplet{
		Id:   s.node.Id,
		Ip:   s.node.Ip,
		Port: s.node.Port,
	}

	s.node.Update(request.Sender)

	return nil
}

func (s Server) FindNode(request *FindRequest, reply *FindNodeResponse) error {
	node := request.Target
	log.Println("FindNode request recieved", node, s)

	shortlist := s.node.ClosestNodes(node.Id, K)
	reply.Triplets = append(reply.Triplets, shortlist...)

	s.node.Update(request.Sender)

	return nil
}

func NewServer(ip string, port string) Server {
	return Server{node: NewNode(ip, port)}
}
