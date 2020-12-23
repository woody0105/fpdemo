package main

import (
	"encoding/binary"
	"encoding/hex"
	// "encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

var clients = make(map[*Client]bool) // connected clients
// var transcribers = make(map[*ffmpeg.Transcriber]bool)

var RandomIDGenerator = func(length uint) string {
	x := make([]byte, length, length)
	for i := 0; i < len(x); i++ {
		x[i] = byte(rand.Uint32())
	}
	return hex.EncodeToString(x)
}

func RandName() string {
	return RandomIDGenerator(10)
}

// func GetTranscriberByID(clientID string) *ffmpeg.Transcriber {
// 	for t := range transcribers {
// 		if t.Id == clientID {
// 			return t
// 		}
// 	}
// 	return nil
// }

func handleconnections1(w http.ResponseWriter, r *http.Request) {
	codec := r.Header.Get("X-WS-Audio-Codec")
	channel := r.Header.Get("X-WS-Audio-Channels")
	sample_rate := r.Header.Get("X-WS-Rate")
	bit_rate := r.Header.Get("X-WS-BitRate")

	if codec == "" || channel == "" || sample_rate == "" || bit_rate == "" {
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
	go handlepkt(w, r, c, codec)
}

func handlepkt(w http.ResponseWriter, r *http.Request, conn *websocket.Conn, codec string) {
	// clientId := RandName()
	// t := ffmpeg.NewTranscriber(clientId)
	// t.Conn = conn
	_, ok := ffmpeg.AudioCodecLookup[codec]
	if !ok {
		panic(fmt.Sprintf("Invalid Codec %s.\n", codec))
	}
	fmt.Println("audio codec id:", codec)
	// var last string
	// var printed string
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// t.TranscriberCodecDeinit()
			// t.StopTranscriber()
			log.Println("read:", err)
			break
		}
		timestamp := binary.BigEndian.Uint64(message[:8])
		packetdata := message[8:]
		timedpacket := ffmpeg.TimedPacket{Timestamp: timestamp, Packetdata: ffmpeg.APacket{packetdata, len(packetdata)}}
		str := ffmpeg.FeedPacket(timedpacket)
		fmt.Println(str)
	}
}

func startServer1() {
	http.HandleFunc("/songdetect", handleconnections1)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	startServer1()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
