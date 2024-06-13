package voiceprocessing

import (
	"errors"

	"github.com/kercre123/0.11-cloud/pkg/speechrequest"
	"github.com/kercre123/0.11-cloud/pkg/stt"
	"github.com/kercre123/0.11-cloud/pkg/ttr"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
)

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
	sReq := speechrequest.ReqToSpeechRequest(req)
	transcribedText, err := stt.STT(sReq)
	if err != nil {
		ttr.IntentPass(req, "intent_system_noaudio", "error: "+err.Error(), map[string]string{}, false)
		return nil, err
	}
	successMatched := ttr.ProcessTextAll(req, transcribedText, LoadedIntents, false)
	if !successMatched {
		ttr.IntentPass(req, "intent_system_unmatched", transcribedText, map[string]string{}, false)
		return nil, errors.New("intent did not match")
	}
	return nil, nil
}
