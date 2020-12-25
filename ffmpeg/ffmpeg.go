package ffmpeg

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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

func ParseRecognitionResult(strIn string) (string, float64, error) {
	strIn = strings.ReplaceAll(strIn, `{'`, `{"`)
	strIn = strings.ReplaceAll(strIn, `, '`, `, "`)
	strIn = strings.ReplaceAll(strIn, `: '`, `: "`)
	strIn = strings.ReplaceAll(strIn, `':`, `":`)
	strIn = strings.ReplaceAll(strIn, `',`, `",`)
	strIn = strings.ReplaceAll(strIn, `'}`, `"}`)
	strIn = strings.ReplaceAll(strIn, `: b'`, `:"`)

	var result map[string]interface{}
	json.Unmarshal([]byte(strIn), &result)
	fmt.Println(result)
	results := result["results"].([]interface{})
	topResult := results[0].(map[string]interface{})
	songName := topResult["song_name"]
	inConfidence := topResult["input_confidence"]
	fmt.Println("input_confidence = ", inConfidence)
	return songName.(string), inConfidence.(float64), nil
}

func Recvpkts2file(recvpkts []TimedPacket) {
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
	res := Recognizefile(workDir + fname)
	fmt.Printf("Recognition result:\n %s", res)
	songname, inConfidence, _ := ParseRecognitionResult(res)
	fmt.Println("songname:", songname, "confidence:", inConfidence)
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
