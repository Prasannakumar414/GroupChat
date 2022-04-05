package main

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *Server) NickEndpoint(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := r.PostForm.Get("name")
	chatroomid := r.PostForm.Get("chatid")

	cookie := &http.Cookie{
		Name:   "nick",
		Value:  name,
		MaxAge: 3000,
	}
	http.SetCookie(w, cookie)

	cookie = &http.Cookie{
		Name:   "chatid",
		Value:  chatroomid,
		MaxAge: 3000,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(200)
}

func (s *Server) WebsocketEndPoint(w http.ResponseWriter, r *http.Request) {
	//get a request check for cookies
	name, chatid, err := GetCookieValues(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, x := s.Chatroommap[chatid]
	if !x {
		s.Chatroommap[chatid] = NewChatroom(chatid)
	}
	chatroom := s.Chatroommap[chatid]
	fmt.Println("check1")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	chatroom.AddUser(name, ws)
	log.Println("Client Successfully Connected!!!")

	results := s.GetLastmsgs(chatid, 10)

	for _, doc := range results {

		Pmsg := fmt.Sprint(doc[3].Value) + " at " + fmt.Sprint(doc[2].Value)
		sender := fmt.Sprint(doc[1].Value)

		if err := SendMessage(sender, Pmsg, ws); err != nil {
			return
		}
	}
	s.ReadWriteMessages(chatroom, name)
}
func (s *Server) setuproutes() {
	fileserver := http.FileServer(http.Dir("/home/prasanna/Programs/Projects/RealtimeChat/frontend"))
	http.Handle("/", fileserver)
	http.HandleFunc("/nick", s.NickEndpoint)
	http.HandleFunc("/ws", s.WebsocketEndPoint)
}

func main() {
	fmt.Println("Server is UP!!!!")
	server := CreateServer()
	defer close(server.client, server.context, server.cancel)
	server.setuproutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
