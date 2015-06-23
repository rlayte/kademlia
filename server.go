package kademlia

import "log"

type Server struct {
	node Node
}

type Request struct {
	Sender *Contact
}

type FindRequest struct {
	Sender *Contact
	Target *Contact
}

type PingResponse struct {
	*Contact
}

type FindNodeResponse struct {
	Contacts []Contact
}

func (s Server) Ping(request *Request, reply *PingResponse) error {
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
	s.node.Update(request.Sender)

	shortlist := s.node.ClosestNodes(node.Id, K)
	reply.Contacts = append(reply.Contacts, shortlist...)

	return nil
}

func NewServer(ip string, port string) Server {
	log.Println("Starting server", ip, port)
	return Server{node: NewNode(ip, port)}
}
