package main

/*
package main

import (
	"fmt"
	"reflect"
	"time"

	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

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

	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)

	var filter, option interface{}

	filter = bson.D{}

	opts := options.Find()
	opts.SetSort(bson.D{{"timestamp", -1}})

	cursor, err := query(client, ctx, "Chatroom",
		"messages", filter, option, 10)
	// handle the errors.
	if err != nil {
		panic(err)
	}

	var results []bson.D

	// to get bson object  from cursor,
	// returns error if any.
	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		log.Println(err)
	}

	// printing the result of query.
	fmt.Println("Query Reult")

	for _, doc := range results {
		pastMessage := doc[3].Value
		sender := doc[1].Value

		Pmsg := fmt.Sprint(pastMessage)
		Send := fmt.Sprint(sender)

		finalMessage := Send + " : " + Pmsg + " " + fmt.Sprint(doc[2].Value)
		fmt.Print(reflect.TypeOf(pastMessage))

		if err := Nicks[name].WriteMessage(1, []byte(finalMessage)); err != nil {
			log.Println(err)
			return
		}

		fmt.Println(pastMessage)

	}

	for {

		var document interface{}

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		document = bson.D{
			{"name", name},
			{"time", time.Now().Format(time.ANSIC)},
			{"message", string(message)},
		}

		insertOneResult, err := insertOne(client, ctx, "Chatroom",
			"messages", document)

		// handle the error
		if err != nil {
			panic(err)
		}

		fmt.Println(insertOneResult.InsertedID)
		nmsg := (name + " : " + string(message) + "  " + time.Now().Format(time.ANSIC))
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

	s := "Hii " + name + " you have entered chatroom " + time.Now().String()

	err = ws.WriteMessage(1, []byte(s))
	if err != nil {
		log.Println(err)
	}

	connectionReader(ws, name)
}

func setuproutes() {
	fileserver := http.FileServer(http.Dir("/home/prasanna/Programs/Projects/RealtimeChat/frontend"))
	http.Handle("/", fileserver)
	http.HandleFunc("/nick", NickEndpoint)
	http.HandleFunc("/ws", websocketpage)
}

func main() {
	fmt.Println("Server is UP")
	setuproutes()
	log.Fatal(http.ListenAndServe(":8080", nil))

	//Chat Room Branch

}
*/
