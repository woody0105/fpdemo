package ffmpeg

import (
	// "sync"
	"unsafe"

	// "github.com/gorilla/websocket"
)

// #cgo pkg-config: libavformat libavfilter libavcodec libavutil libswscale gnutls
// #include <stdlib.h>
// #include "lpms_ffmpeg.h"
import "C"

const (
	MP2 = iota + 0x15000
	MP3
	AAC
)

var AudioCodecLookup = map[string]int{
	"AAC": AAC,
	"MP3": MP3,
	"MP2": MP2,
}

type TimedPacket struct {
	Packetdata APacket
	Timestamp  uint64
}

type APacket struct {
	Data   []byte
	Length int
}

func FeedPacket(pkt TimedPacket) string {
	pktdata := pkt.Packetdata
	buffer := (*C.char)(unsafe.Pointer(C.CString(string(pktdata.Data))))
	defer C.free(unsafe.Pointer(buffer))
	str := C.ds_feedpkt(buffer, C.int(pktdata.Length))
	return C.GoString(str)
}

func CodecInit() {
	C.audio_codec_init()
}

func CodecDeinit() {
	C.audio_codec_deinit()
}

// func NewTranscriber(Id string) *Transcriber {
// 	t := &Transcriber{
// 		Id:           Id,
// 		codec_params: C.lpms_codec_new(),
// 		mu:           &sync.Mutex{},
// 		streamState:  C.t_create_stream(),
// 		refeed_data:  C.t_refeed_data(),
// 	}
// 	log.Println("New transcriber created.")
// 	return t
// }

// func (t *Transcriber) TranscriberCodecInit(codec_id int) {
// 	codec_params := t.codec_params
// 	C.t_audio_codec_init(C.int(codec_id), codec_params)
// }

// func (t *Transcriber) TranscriberCodecDeinit() {
// 	codec_params := t.codec_params
// 	C.t_audio_codec_deinit(codec_params)
// }

// func (t *Transcriber) FeedPacket(pkt TimedPacket) string {
// 	pktdata := pkt.Packetdata
// 	codec_params := t.codec_params
// 	stream_ctx := t.streamState
// 	buffer := (*C.char)(unsafe.Pointer(C.CString(string(pktdata.Data))))
// 	defer C.free(unsafe.Pointer(buffer))
// 	str := (*C.char)(unsafe.Pointer(C.malloc(C.sizeof_char * 256)))
// 	defer C.free(unsafe.Pointer(str))
// 	refeed_data := t.refeed_data

// 	t.mu.Lock()
// 	new_stream_ctx := C.t_ds_feedpkt(codec_params, stream_ctx, buffer, C.int(pktdata.Length), refeed_data, str)
// 	t.refeed_data = refeed_data
// 	if new_stream_ctx != nil {
// 		t.streamState = new_stream_ctx
// 	}
// 	t.mu.Unlock()
// 	return C.GoString(str)
// }
