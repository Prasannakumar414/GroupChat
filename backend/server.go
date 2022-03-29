package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Nicks map[string]*websocket.Conn

func init() {
	Nicks = make(map[string]*websocket.Conn)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func connectionReader(conn *websocket.Conn, name string) {

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		nmsg := (name + " : " + string(message))
		for key, val := range Nicks {
			if key != name {
				if err := val.WriteMessage(messageType, []byte(nmsg)); err != nil {
					log.Println(err)
					return
				}
			}
		}
		fmt.Println("Sent!!")

	}
}

func NickEndpoint(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	name := string(r.PostForm.Get("name"))

	fmt.Println(name)
	cookie := &http.Cookie{
		Name:   "nick",
		Value:  name,
		MaxAge: 3000,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(200)
}

func websocketpage(w http.ResponseWriter, r *http.Request) {

	tokenCookie, err := r.Cookie("nick")
	if err != nil {
		log.Println("Error While Reading Cookie")
		w.WriteHeader(403)
		return
	}
	fmt.Println("\nPrinting cookie with name as token")
	name := string(tokenCookie.Value)
	fmt.Println(name)
	fmt.Println()

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	Nicks[name] = ws
	log.Println("Client Successfully Connected")

	s := "Hii " + name + " you have entered chatroom"

	err = ws.WriteMessage(1, []byte(s))
	if err != nil {
		log.Println(err)
	}

	connectionReader(ws, name)
}

func setuproutes() {
	fileserver := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fileserver)
	http.HandleFunc("/nick", NickEndpoint)
	http.HandleFunc("/ws", websocketpage)
}

func main() {
	fmt.Println("Server is UP")
	setuproutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
