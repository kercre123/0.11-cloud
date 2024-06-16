package voiceprocessing

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	sr "github.com/kercre123/0.11-cloud/pkg/speechrequest"
	"github.com/soundhound/houndify-sdk-go"
)

type inMemorySeeker struct {
	buf *bytes.Buffer
	pos int64
}

func (ims *inMemorySeeker) Write(p []byte) (n int, err error) {
	n, err = ims.buf.Write(p)
	ims.pos += int64(n)
	return n, err
}

func (ims *inMemorySeeker) Seek(offset int64, whence int) (int64, error) {
	var newPos int64
	switch whence {
	case io.SeekStart:
		newPos = offset
	case io.SeekCurrent:
		newPos = ims.pos + offset
	case io.SeekEnd:
		newPos = int64(ims.buf.Len()) + offset
	default:
		return 0, fmt.Errorf("invalid whence: %d", whence)
	}
	if newPos < 0 {
		return 0, fmt.Errorf("negative position")
	}
	ims.pos = newPos
	return ims.pos, nil
}

func (ims *inMemorySeeker) Bytes() []byte {
	return ims.buf.Bytes()
}

func WAVEncode(pcmData []byte) ([]byte, error) {
	buf := &inMemorySeeker{buf: new(bytes.Buffer)}
	enc := wav.NewEncoder(buf, 16000, 16, 1, 1)
	audioBuffer := &audio.IntBuffer{
		Data: make([]int, len(pcmData)/2),
		Format: &audio.Format{
			SampleRate:  16000,
			NumChannels: 1,
		},
	}
	for i := 0; i < len(pcmData); i += 2 {
		audioBuffer.Data[i/2] = int(int16(pcmData[i]) | int16(pcmData[i+1])<<8)
	}
	if err := enc.Write(audioBuffer); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func StreamAudioToHoundify(sreq sr.SpeechRequest, client houndify.Client) string {
	var err error
	rp, wp := io.Pipe()
	req := houndify.VoiceRequest{
		AudioStream: rp,
		UserID:      sreq.Device,
		RequestID:   sreq.Session,
	}
	done := make(chan bool)
	//speechDone := false
	go func(wp *io.PipeWriter) {
		defer wp.Close()

		for {
			select {
			case <-done:
				return
			default:
				var chunk []byte
				chunk, err = sreq.GetNextStreamChunkOpus()
				//speechDone, _ = sreq.DetectEndOfSpeech()
				if err != nil {
					fmt.Println("End of stream")
					return
				}

				if sreq.IsOpus {
					wp.Write(chunk)
				} else {
					encodedChunk, err := WAVEncode(chunk)
					if err != nil {
						fmt.Println(err)
					}
					wp.Write(encodedChunk)
				}

				// if speechDone {
				// 	return
				// }
			}
		}
	}(wp)

	partialTranscripts := make(chan houndify.PartialTranscript)
	go func() {
		for partial := range partialTranscripts {
			if *partial.SafeToStopAudio {
				fmt.Println("SafeToStopAudio received")
				done <- true
				return
			}
		}
	}()

	serverResponse, err := client.VoiceSearch(req, partialTranscripts)
	if err != nil {
		fmt.Println(err)
		fmt.Println(serverResponse)
	}
	return serverResponse
}
