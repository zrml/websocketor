package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type book struct {
	ISBN   int    //`json:"isbn"`
	Title  string //`json:"title"`
	Author string //`json:"author"`
}

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {

	log.Println("URL/ws endpoint called and followed...")

	// and for now allow any connection from any origin (no checks)
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// checks...
	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	wsReader(ws)

	defer ws.Close()
}

// ws reader: will continually listen for new messages sent to our WebSocket endpoint
func wsReader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// print out that message for clarity
		log.Println("Received from client: ", string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

		var receivedString = string(p)
		log.Printf("receivedString=%s", receivedString)

		// miminking the call to a xDBC driver for IRIS Cloud SQL
		if receivedString == "SELECT *" {
			log.Printf("Received=%s", receivedString)

			// fetch DB data
			bookData, err := fetchData(receivedString)
			if err != nil {
				log.Printf("ERROR calling fetchData(); %s", err)
			}

			//log.Printf("Returned bookData=%s", string(bookData))

			// write the JSON result-set back to the client ws listener
			//err = conn.WriteMessage(messageType, []byte("QUERY_Returned_DATA..."))
			err = conn.WriteMessage(messageType, bookData)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func fetchData(qString string) (tableData []byte, err error) {

	if qString == "" {
		log.Println("qString is empty")

	} else {
		bData := []byte(`
			[
			{"ISBN": 111, "Title": "The DeBono Code", "Author": "Edward DeBono" },
			{"ISBN": 222, "Title": "The Republic", "Author": "Plato"},
			{"ISBN": 333, "Title": "La Divina Commedia", "Author": "Dante"}
			]`)

		// assignment for returning it (Go var scope!)
		tableData = bData

		log.Printf("JSON: %s", bData)

		// unmashalling for log checking...
		/*var bks []book
		err := json.Unmarshal(bData, &bks)
		if err != nil {
			log.Println(err)
		}

		for i := range bks {
			//fmt.Printf("ISBN=%d; Title=%s; Author=%s\n", bks[i].ISBN, bks[i].Title, bks[i].Author)
			log.Printf("ISBN=%d; Title=%s; Author=%s", bks[i].ISBN, bks[i].Title, bks[i].Author)
		}
		*/
	}

	return
}

func homePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Default http page called")
	fmt.Fprintf(w, "Default Home Page")
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	//fmt.Printf("\nHello Websocketor! You are switched ON.) \n")
	log.Println("The Websocketor is ON")

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
