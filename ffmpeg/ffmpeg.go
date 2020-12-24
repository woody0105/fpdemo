package ffmpeg

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
