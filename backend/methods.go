package main

import (
	//"context"

	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"

	//"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateServer() *Server {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)

	return NewServer(client, ctx)
}

func (s *Server) GetLastmsgs(chatroomid string, limit int64) []bson.D {
	var filter, option interface{}
	filter = bson.D{}

	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})

	cursor, err := query(s.client, s.context, chatroomid,
		"messages", filter, option, limit)
	if err != nil {
		log.Println(err)
	}

	var results []bson.D
	if err := cursor.All(s.context, &results); err != nil {
		log.Println(err)
	}
	return results
}

func (s *Server) StoreMyMessage(message string, chatroomid string, name string) {
	document := bson.D{
		{"name", name},
		{"time", time.Now().Format(time.ANSIC)},
		{"message", string(message)},
	}

	insertOneResult, err := insertOne(s.client, s.context, chatroomid,
		"messages", document)
	if err != nil {
		panic(err)
	}
	fmt.Println(insertOneResult.InsertedID)
}

func (c *Chatroom) Sendmsg(msg string, name string) {
	for key, val := range c.Nicks {
		if key != name {
			if err := SendMessage(name, msg, val); err != nil {
				return
			}
		}
	}
}

func GetCookieValues(r *http.Request) (string, string, error) {
	tokenCookie, err1 := r.Cookie("nick")
	chatCookie, err2 := r.Cookie("chatid")
	if err1 != nil {
		log.Println("Error While Reading Cookie")
		return "", "", err1
	}
	if err2 != nil {
		log.Println("Error While Reading Cookie")
		return "", "", err2
	}
	name := string(tokenCookie.Value)
	chatid := string(chatCookie.Value)

	return name, chatid, nil
}

func SendMessage(name string, msg string, ws *websocket.Conn) error {
	nmsg := (name + " : " + string(msg) + "  " + time.Now().Format(time.ANSIC))
	err := ws.WriteMessage(1, []byte(nmsg))
	if err != nil {
		log.Println(err)
	}
	return err
}
