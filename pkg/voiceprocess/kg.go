package voiceprocessing

import (
	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
)

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	// sReq := speechrequest.ReqToSpeechRequest(req)
	// transcribedText, err := stt.STT(sReq)
	// if err != nil {
	// 	ttr.IntentPass(req, "intent_system_noaudio", "error: "+err.Error(), map[string]string{}, false)
	// 	return nil, err
	// }
	// successMatched := ttr.ProcessTextAll(req, transcribedText, LoadedIntents, false)
	// if !successMatched {
	// 	ttr.IntentPass(req, "intent_system_unmatched", transcribedText, map[string]string{}, false)
	// 	return nil, errors.New("intent did not match")
	// }
	req.Stream.Send(&pb.KnowledgeGraphResponse{
		QueryText:   "unknown",
		SpokenText:  "Knowledge Graph is not implemented yet.",
		CommandType: "unknown",
		DomainsUsed: []string{"unknown"},
	})
	return nil, nil
}
