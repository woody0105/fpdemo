package ffmpeg

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
)

type TimedPacket struct {
	Packetdata APacket
	Timestamp  uint64
}

type APacket struct {
	Data   []byte
	Length int
}

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

func ParseRecognitionResult(strIn string) (string, string, string, float64, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()
	strIn = strings.ReplaceAll(strIn, `{'`, `{"`)
	strIn = strings.ReplaceAll(strIn, `, '`, `, "`)
	strIn = strings.ReplaceAll(strIn, `: '`, `: "`)
	strIn = strings.ReplaceAll(strIn, `':`, `":`)
	strIn = strings.ReplaceAll(strIn, `',`, `",`)
	strIn = strings.ReplaceAll(strIn, `'}`, `"}`)
	strIn = strings.ReplaceAll(strIn, `: b'`, `:"`)

	var result map[string]interface{}
	json.Unmarshal([]byte(strIn), &result)
	results := result["results"].([]interface{})
	topResult := results[0].(map[string]interface{})
	inConfidence := topResult["input_confidence"]
	songTitle := topResult["song_title"]
	artist := topResult["artist"]
	songName := topResult["song_name"]
	return songName.(string), songTitle.(string), artist.(string), inConfidence.(float64), nil
}

// func Recvpkts2file(recvpkts []TimedPacket) {
// 	fname := "clip" + RandName() + ".aac"
// 	workDir := ".tmp/"
// 	file, err := os.Create(workDir + fname)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, pkt := range recvpkts {
// 		pktdata := pkt.Packetdata
// 		file.Write(pktdata.Data)
// 	}
// 	file.Close()
// 	res := Recognizefile(workDir + fname)
// 	fmt.Printf("Recognition result:\n %s", res)
// 	songname, inConfidence, _ := ParseRecognitionResult(res)
// 	fmt.Println("songname:", songname, "confidence:", inConfidence)
// }

func ProcessPkt(recvpkts []TimedPacket, conn *websocket.Conn) {
	fname := "clip" + RandName() + ".aac"
	workDir := ".tmp/"
	file, err := os.Create(workDir + fname)
	if err != nil {
		panic(err)
	}

	for _, pkt := range recvpkts {
		pktdata := pkt.Packetdata
		file.Write(pktdata.Data)
	}
	file.Close()
	recogRes := Recognizefile(workDir + fname)
	songName, songTitle, artist, inConfidence, err := ParseRecognitionResult(recogRes)

	if err != nil {
		fmt.Println("Recognition result parsing failed.")
		return
	}

	if inConfidence <= 0.07 {
		fmt.Println("no matching result")
		return
	}

	res := map[string]interface{}{"songname": songName, "songtitle": songTitle, "artist": artist}
	jsonres, err := json.Marshal(res)
	fmt.Println(string(jsonres))
	conn.WriteMessage(websocket.TextMessage, []byte(string(jsonres)))
}

func Recognizefile(fname string) string {
	args := []string{
		"dejavu.py", "-r",
		"file", fname,
	}
	out, _ := exec.Command("python3", args...).Output()
	os.Remove(fname)
	return string(out)
}
