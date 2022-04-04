package main

import (
	"context"

	"github.com/gorilla/websocket"

	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Chatroom struct {
	Chatroomid string
	Nicks      map[string]*websocket.Conn
}
type Server struct {
	Chatroommap map[string]*Chatroom
	client      *mongo.Client
	context     context.Context
}

func NewChatroom(chatroomid string) *Chatroom {
	return &Chatroom{
		Chatroomid: chatroomid,
		Nicks:      make(map[string]*websocket.Conn),
	}
}

func NewServer(clientname *mongo.Client, context context.Context) *Server {
	return &Server{
		Chatroommap: make(map[string]*Chatroom),
		client:      clientname,
		context:     context,
	}
}

func (s *Server) AddChatroom(chatroomid string) *Chatroom {
	x := NewChatroom(chatroomid)
	s.Chatroommap[chatroomid] = x
	return x
}
func (c *Chatroom) AddUser(name string, WSconnection *websocket.Conn) {
	c.Nicks[name] = WSconnection
}
