package kademlia

import "log"

type Server struct {
	node Node
}

type Request struct {
	Sender *Contact
}

type FindRequest struct {
	*Request
	Target *Contact
}

type PingResponse struct {
	*Contact
}

type FindNodeResponse struct {
	Contacts []Contact
}

func addToShortlist(shortlist chan Contact, items []Contact) {
	for _, contact := range items {
		shortlist <- contact
	}
}

func requestClosest(contact Contact, method string) []Contact {
	reply := FindNodeResponse{}
	contact.client.Call(method, &contact, reply)

	return reply.Contacts
}

func iterateShortlist(shortlist chan Contact, target NodeId, method string) Contact {
	var closestNode Contact
	var closestDistance NodeId
	done := make(chan Contact, 1)
	contacted := map[NodeId]bool{}

	for contact := range shortlist {
		if len(contacted) > K {
			close(shortlist)
			done <- closestNode
		}

		if distance := Xor(target, contact.Id); distance.LessThan(closestDistance) {
			closestDistance = distance
			closestNode = contact
		}

		if _, ok := contacted[contact.Id]; !ok {
			contacted[contact.Id] = true
			addToShortlist(shortlist, requestClosest(contact, method))
		}
	}

	return <-done
}

func nodeLookup(origin Node, target NodeId, method string) Contact {
	shortlist := make(chan Contact, A)
	addToShortlist(shortlist, origin.ClosestNodes(target, A))
	return iterateShortlist(shortlist, target, method)
}

func (s Server) Ping(request *Request, reply *PingResponse) error {
	log.Println("Ping recieved:", request.Sender)

	reply.Contact = &Contact{
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
	reply.Contacts = append(reply.Contacts, shortlist...)

	s.node.Update(request.Sender)

	return nil
}

func NewServer(ip string, port string) Server {
	return Server{node: NewNode(ip, port)}
}
