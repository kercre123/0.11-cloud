package voiceprocessing

import (
	"errors"
	"net"

	"github.com/kercre123/0.11-cloud/pkg/speechrequest"
	"github.com/kercre123/0.11-cloud/pkg/stt"
	"github.com/kercre123/0.11-cloud/pkg/ttr"
	"github.com/kercre123/0.11-cloud/pkg/vars"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
	"google.golang.org/grpc/peer"
)

func (s *Server) ProcessIntent(req *vtt.IntentRequest) (*vtt.IntentResponse, error) {
	p, _ := peer.FromContext(req.Stream.Context())
	ip, _, _ := net.SplitHostPort(p.Addr.String())
	vars.AddToIPWhitelist(ip, req.Device)
	sReq := speechrequest.ReqToSpeechRequest(req)
	transcribedText, err := stt.STT(sReq, true)
	if err != nil {
		ttr.IntentPass(req, "intent_system_noaudio", "error: "+err.Error(), map[string]string{}, false)
		return nil, err
	}
	successMatched := ttr.ProcessTextAll(req, transcribedText, LoadedIntents, sReq.IsOpus)
	if !successMatched {
		ttr.IntentPass(req, "intent_system_unmatched", transcribedText, map[string]string{}, false)
		return nil, errors.New("intent did not match")
	}
	return nil, nil
}
