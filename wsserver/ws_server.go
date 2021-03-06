package main

import (
	"encoding/binary"

	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	// "strings"

	"github.com/gorilla/websocket"
	"github.com/woody0105/fpdemo/ffmpeg"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Client struct {
	ID   string
	Conn *websocket.Conn
}

type message struct {
	clientid string
	data     []byte
}

var msgchan = make(chan message)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleconnections1(w http.ResponseWriter, r *http.Request) {
	codec := r.Header.Get("X-WS-Audio-Codec")
	channel := r.Header.Get("X-WS-Audio-Channels")
	sampleRate := r.Header.Get("X-WS-Rate")
	bitRate := r.Header.Get("X-WS-BitRate")

	if codec == "" || channel == "" || sampleRate == "" || bitRate == "" {
		log.Print("audio meta data not present in header, handshake failed.")
		return
	}
	respheader := make(http.Header)
	respheader.Add("Sec-WebSocket-Protocol", "speechtotext.livepeer.com")
	c, err := upgrader.Upgrade(w, r, respheader)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go handlepkt(w, r, c)
}

func handlepkt(w http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	var recvpkt []ffmpeg.TimedPacket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		timestamp := binary.BigEndian.Uint64(message[:8])
		packetdata := message[8:]
		timedpacket := ffmpeg.TimedPacket{Timestamp: timestamp, Packetdata: ffmpeg.APacket{packetdata, len(packetdata)}}
		recvpkt = append(recvpkt, timedpacket)
		if len(recvpkt) > 700 {
			fmt.Printf("Processing packets...\n")
			// go ffmpeg.Recvpkts2file(recvpkt)
			go ffmpeg.ProcessPkt(recvpkt, conn)
			recvpkt = nil
		}
	}
}

func startServer1() {
	http.HandleFunc("/songdetect", handleconnections1)
	log.Print("Listening at /songdetect")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	startServer1()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
