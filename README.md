# Go websocket example

## What is it?
It's a simple demo of a 
- server side websocket implemented in Go and
- the client side web-page with the relevant JS

## How to run it
1) run the server 
	- go run main.go
2) make sure the nodejs webgateway _live-server_, is installed 
	- npm install -g live-server
	- it will throw you into a webpage with all files listed in the project
	- select/click "page-with-ws-call.html"
3) run the exercises in the web page
	- in the input field, try any string and also "SELECT *"


## Interesting points

### Backend
I'm using the Gorilla package.<br>
One must upgrade the protocol from http to websockets.<br>
See the Upgrader.<br>


### Frontend
In the 3 main parts of the page it shows that
- it calls a default http-based page
- it sends an automated message (open developer/console in the browser)
- it sends and echos from the server a string and the status
- it retuns some data if you type "SELECT *"

