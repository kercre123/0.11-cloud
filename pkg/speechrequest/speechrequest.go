package speechrequest

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/opus-go/opus"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
	"github.com/maxhawkins/go-webrtcvad"
)

// one type and many functions for dealing with intent, intent-graph, and knowledge-graph requests
// also some functions to help decode the stream bytes into ones friendly for stt engines

var debugWriteFile bool = false
var debugFile *os.File

type SpeechRequest struct {
	Device         string
	Session        string
	FirstReq       []byte
	Stream         interface{}
	IsKG           bool
	IsIG           bool
	MicData        []byte
	DecodedMicData []byte
	// for PCM bots
	EncodedMicData []byte
	PrevLen        int
	PrevLenRaw     int
	InactiveFrames int
	ActiveFrames   int
	VADInst        *webrtcvad.VAD
	LastAudioChunk []byte
	IsOpus         bool
	OpusStream     *opus.OggStream
	OpusStreamEnc  *opus.OggStream
	EncBuf         []byte
	ChunkSkips     int
}

func BytesToSamples(buf []byte) []int16 {
	samples := make([]int16, len(buf)/2)
	for i := 0; i < len(buf)/2; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(buf[i*2:]))
	}
	return samples
}

func (req *SpeechRequest) OpusDetect() bool {
	var isOpus bool
	if len(req.FirstReq) > 0 {
		if req.FirstReq[0] == 0x4f {
			fmt.Println("Bot " + req.Device + " Stream type: OPUS")
			isOpus = true
		} else {
			isOpus = false
			fmt.Println("Bot " + req.Device + " Stream type: PCM")
		}
	}
	return isOpus
}

func (req *SpeechRequest) OpusDecode(chunk []byte) []byte {
	if req.IsOpus {
		n, err := req.OpusStream.Decode(chunk)
		if err != nil {
			fmt.Println(err)
		}
		return n
	} else {
		return chunk
	}
}

func SplitVAD(buf []byte) [][]byte {
	var chunk [][]byte
	for len(buf) >= 320 {
		chunk = append(chunk, buf[:320])
		buf = buf[320:]
	}
	return chunk
}

func BytesToIntVAD(stream opus.OggStream, data []byte, die bool, isOpus bool) [][]byte {
	// detect if data is pcm or opus
	if die {
		return nil
	}
	if isOpus {
		// opus
		n, err := stream.Decode(data)
		if err != nil {
			fmt.Println(err)
		}
		byteArray := SplitVAD(n)
		return byteArray
	} else {
		// pcm
		byteArray := SplitVAD(data)
		return byteArray
	}
}

// Uses VAD to detect when the user stops speaking
func (req *SpeechRequest) DetectEndOfSpeech() (bool, bool) {
	// changes InactiveFrames and ActiveFrames in req
	inactiveNumMax := 23
	vad := req.VADInst
	vad.SetMode(2)
	for _, chunk := range SplitVAD(req.LastAudioChunk) {
		active, err := vad.Process(16000, chunk)
		if err != nil {
			fmt.Println("VAD err:")
			fmt.Println(err)
			return true, false
		}
		if active {
			req.ActiveFrames = req.ActiveFrames + 1
			req.InactiveFrames = 0
		} else {
			req.InactiveFrames = req.InactiveFrames + 1
		}
		if req.InactiveFrames >= inactiveNumMax && req.ActiveFrames > 30 {
			fmt.Println("(Bot " + req.Device + ") End of speech detected.")
			return true, true
		}
	}
	if req.ActiveFrames < 5 {
		return false, false
	}
	return false, true
}

// Converts a vtt.*Request to a SpeechRequest, which allows functions like DetectEndOfSpeech to work
func ReqToSpeechRequest(req interface{}) SpeechRequest {
	if debugWriteFile {
		debugFile, _ = os.Create("/tmp/wirepodtest.ogg")
	}
	var request SpeechRequest
	request.PrevLen = 0
	var err error
	request.VADInst, err = webrtcvad.New()
	if err != nil {
		fmt.Println(err)
	}
	if str, ok := req.(*vtt.IntentRequest); ok {
		var req1 *vtt.IntentRequest = str
		request.Device = req1.Device
		request.Session = req1.Session
		request.Stream = req1.Stream
		request.FirstReq = req1.FirstReq.InputAudio
		request.MicData = append(request.MicData, req1.FirstReq.InputAudio...)
	} else if str, ok := req.(*vtt.KnowledgeGraphRequest); ok {
		var req1 *vtt.KnowledgeGraphRequest = str
		request.IsKG = true
		request.Device = req1.Device
		request.Session = req1.Session
		request.Stream = req1.Stream
		request.FirstReq = req1.FirstReq.InputAudio
		request.MicData = append(request.MicData, req1.FirstReq.InputAudio...)
	} else if str, ok := req.(*vtt.IntentGraphRequest); ok {
		request.IsIG = true
		var req1 *vtt.IntentGraphRequest = str
		request.Device = req1.Device
		request.Session = req1.Session
		request.Stream = req1.Stream
		request.FirstReq = req1.FirstReq.InputAudio
		if debugWriteFile {
			debugFile.Write(req1.FirstReq.InputAudio)
		}
		request.MicData = append(request.MicData, req1.FirstReq.InputAudio...)
	} else {
		fmt.Println("reqToSpeechRequest: invalid type")
	}
	isOpus := request.OpusDetect()
	if isOpus {
		request.OpusStream = &opus.OggStream{}
		decodedFirstReq, _ := request.OpusStream.Decode(request.FirstReq)
		request.FirstReq = decodedFirstReq
		request.DecodedMicData = append(request.DecodedMicData, decodedFirstReq...)
		request.LastAudioChunk = request.DecodedMicData[request.PrevLen:]
		request.PrevLen = len(request.DecodedMicData)
		request.IsOpus = true
	}
	//  else {
	// 	request.OpusStreamEnc = &opus.OggStream{
	// 		SampleRate: 16000,
	// 		Channels:   1,
	// 	}
	// 	encodedFirstReq, err := request.OpusStreamEnc.EncodeBytes(request.FirstReq)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	request.EncodedMicData = append(request.DecodedMicData, encodedFirstReq...)
	// }
	return request
}

// Returns the next chunk in the stream as 16000 Hz PCM
func (req *SpeechRequest) GetNextStreamChunk() ([]byte, error) {
	// returns next chunk in voice stream as pcm
	if str, ok := req.Stream.(pb.ChipperGrpc_StreamingIntentServer); ok {
		var stream pb.ChipperGrpc_StreamingIntentServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.DecodedMicData[req.PrevLen:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		return dataReturn, nil
	} else if str, ok := req.Stream.(pb.ChipperGrpc_StreamingIntentGraphServer); ok {
		var stream pb.ChipperGrpc_StreamingIntentGraphServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.DecodedMicData[req.PrevLen:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		if debugWriteFile {
			debugFile.Write(chunk.InputAudio)
		}
		return dataReturn, nil
	} else if str, ok := req.Stream.(pb.ChipperGrpc_StreamingKnowledgeGraphServer); ok {
		var stream pb.ChipperGrpc_StreamingKnowledgeGraphServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.DecodedMicData[req.PrevLen:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		return dataReturn, nil
	}
	fmt.Println("invalid type")
	return nil, errors.New("invalid type")
}

// Returns next chunk in the stream as whatever the original format is (OPUS 99% of the time)
func (req *SpeechRequest) GetNextStreamChunkOpus() ([]byte, error) {
	if str, ok := req.Stream.(pb.ChipperGrpc_StreamingIntentServer); ok {
		var stream pb.ChipperGrpc_StreamingIntentServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.MicData[req.PrevLenRaw:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		req.PrevLenRaw = len(req.MicData)
		return dataReturn, nil
	} else if str, ok := req.Stream.(pb.ChipperGrpc_StreamingIntentGraphServer); ok {
		var stream pb.ChipperGrpc_StreamingIntentGraphServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.MicData[req.PrevLenRaw:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		req.PrevLenRaw = len(req.MicData)
		return dataReturn, nil
	} else if str, ok := req.Stream.(pb.ChipperGrpc_StreamingKnowledgeGraphServer); ok {
		var stream pb.ChipperGrpc_StreamingKnowledgeGraphServer = str
		chunk, chunkErr := stream.Recv()
		if chunkErr != nil {
			fmt.Println(chunkErr)
			return nil, chunkErr
		}
		// for Houndify, we must re-encode to Opus
		// if !req.IsOpus {
		// 	req.MicData = append(req.MicData, chunk.InputAudio...)
		// 	encoded, err := req.OpusStreamEnc.EncodeBytes(chunk.InputAudio)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// 	req.EncodedMicData = append(req.EncodedMicData, encoded...)
		// 	req.LastAudioChunk = req.EncodedMicData[req.PrevLen:]
		// 	req.PrevLen = len(req.EncodedMicData)
		// 	req.PrevLenRaw = len(req.MicData)
		// 	return req.LastAudioChunk, nil
		// }
		req.MicData = append(req.MicData, chunk.InputAudio...)
		req.DecodedMicData = append(req.DecodedMicData, req.OpusDecode(chunk.InputAudio)...)
		dataReturn := req.MicData[req.PrevLenRaw:]
		req.LastAudioChunk = req.DecodedMicData[req.PrevLen:]
		req.PrevLen = len(req.DecodedMicData)
		req.PrevLenRaw = len(req.MicData)
		return dataReturn, nil
	}
	fmt.Println("invalid type")
	return nil, errors.New("invalid type")
}
