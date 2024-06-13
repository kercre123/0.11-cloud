package ttr

import (
	"fmt"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/kercre123/0.11-cloud/pkg/vars"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
)

func IntentPass(req interface{}, intentThing string, speechText string, intentParams map[string]string, isParam bool) (interface{}, error) {
	var esn string
	var req1 *vtt.IntentRequest
	var req2 *vtt.IntentGraphRequest
	var isIntentGraph bool
	if str, ok := req.(*vtt.IntentRequest); ok {
		req1 = str
		esn = req1.Device
		isIntentGraph = false
	} else if str, ok := req.(*vtt.IntentGraphRequest); ok {
		req2 = str
		esn = req2.Device
		isIntentGraph = true
	}

	var intentResult pb.IntentResult
	if isParam {
		intentResult = pb.IntentResult{
			QueryText:  speechText,
			Action:     intentThing,
			Parameters: intentParams,
		}
	} else {
		intentResult = pb.IntentResult{
			QueryText: speechText,
			Action:    intentThing,
		}
	}
	intent := pb.IntentResponse{
		IsFinal:      true,
		IntentResult: &intentResult,
	}
	intentGraphSend := pb.IntentGraphResponse{
		ResponseType: pb.IntentGraphMode_INTENT,
		IsFinal:      true,
		IntentResult: &intentResult,
		CommandType:  pb.RobotMode_VOICE_COMMAND.String(),
	}
	if !isIntentGraph {
		if err := req1.Stream.Send(&intent); err != nil {
			return nil, err
		}
		r := &vtt.IntentResponse{
			Intent: &intent,
		}
		fmt.Println("Bot " + esn + " Intent Sent: " + intentThing)
		if isParam {
			fmt.Println("Bot "+esn+" Parameters Sent:", intentParams)
		} else {
			fmt.Println("No Parameters Sent")
		}
		return r, nil
	} else {
		if err := req2.Stream.Send(&intentGraphSend); err != nil {
			return nil, err
		}
		r := &vtt.IntentGraphResponse{
			Intent: &intentGraphSend,
		}
		fmt.Println("Bot " + esn + " Intent Sent: " + intentThing)
		if isParam {
			fmt.Println("Bot "+esn+" Parameters Sent:", intentParams)
		} else {
			fmt.Println("No Parameters Sent")
		}
		return r, nil
	}
}

func ProcessTextAll(req interface{}, voiceText string, intents []vars.JsonIntent, isOpus bool) bool {
	var botSerial string
	var req2 *vtt.IntentRequest
	var req1 *vtt.KnowledgeGraphRequest
	var req3 *vtt.IntentGraphRequest
	if str, ok := req.(*vtt.IntentRequest); ok {
		req2 = str
		botSerial = req2.Device
	} else if str, ok := req.(*vtt.KnowledgeGraphRequest); ok {
		req1 = str
		botSerial = req1.Device
	} else if str, ok := req.(*vtt.IntentGraphRequest); ok {
		req3 = str
		botSerial = req3.Device
	}
	var matched int = 0
	var intentNum int = 0
	var successMatched bool = false
	voiceText = strings.ToLower(voiceText)
	fmt.Println("Not a custom intent")
	// Look for a perfect match first
	for _, b := range intents {
		for _, c := range b.Keyphrases {
			if voiceText == strings.ToLower(c) {
				fmt.Println("Bot " + botSerial + " Perfect match for intent " + b.Name + " (" + strings.ToLower(c) + ")")
				prehistoricParamChecker(req, b.Name, voiceText, botSerial)
				successMatched = true
				matched = 1
				break
			}
		}
		if matched == 1 {
			matched = 0
			break
		}
		intentNum = intentNum + 1
	}
	// Not found? Then let's be happy with a bare substring search
	if !successMatched {
		intentNum = 0
		matched = 0
		for _, b := range intents {
			for _, c := range b.Keyphrases {
				if strings.Contains(voiceText, strings.ToLower(c)) && !b.RequireExactMatch {
					fmt.Println("Bot " + botSerial + " Partial match for intent " + b.Name + " (" + strings.ToLower(c) + ")")
					prehistoricParamChecker(req, b.Name, voiceText, botSerial)
					successMatched = true
					matched = 1
					break
				}
			}
			if matched == 1 {
				matched = 0
				break
			}
			intentNum = intentNum + 1
		}
	}
	return successMatched
}
