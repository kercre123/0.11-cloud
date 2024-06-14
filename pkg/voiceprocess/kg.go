package voiceprocessing

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/kercre123/0.11-cloud/pkg/speechrequest"
	"github.com/kercre123/0.11-cloud/pkg/stt"
	"github.com/kercre123/0.11-cloud/pkg/vars"
	"github.com/kercre123/0.11-cloud/pkg/vtt"
	"github.com/sashabaranov/go-openai"
	"github.com/soundhound/houndify-sdk-go"
)

func ParseSpokenResponse(serverResponseJSON string) (string, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(serverResponseJSON), &result)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("failed to decode json")
	}
	if !strings.EqualFold(result["Status"].(string), "OK") {
		return "", errors.New(result["ErrorMessage"].(string))
	}
	if result["NumToReturn"].(float64) < 1 {
		return "", errors.New("no results to return")
	}
	return result["AllResults"].([]interface{})[0].(map[string]interface{})["SpokenResponseLong"].(string), nil
}

func removeSpecialCharacters(str string) string {
	re := regexp.MustCompile(`[&^*#@]`)
	return re.ReplaceAllString(str, "")
}

func (s *Server) ProcessKnowledgeGraph(req *vtt.KnowledgeGraphRequest) (*vtt.KnowledgeGraphResponse, error) {
	var spokenText string
	var defaultllmPrompt string = "You are a helpful, animated robot called Vector. Keep the response concise yet informative."
	var llmPromptAddition string = "\n\n" + "The text to speech utterance can only be 255 characters long. BE CONCISE. The user input might not be spelt/punctuated correctly as it is coming from speech-to-text software. Do not include special characters in your answer. This includes the following characters (not including the quotes): '& ^ * # @ -'. DON'T INCLUDE THESE. DON'T MAKE LISTS WITH FORMATTING. THINK OF THE SPEECH-TO-TEXT ENGINE. If you want to use a hyphen, Use it like this: 'something something -- something -- something something'."
	sReq := speechrequest.ReqToSpeechRequest(req)
	_, uInfo := vars.GetUserInfo(req.Device)
	if !uInfo.KG.Enabled {
		spokenText = "Knowledge Graph is not enabled for this bot."
	} else if uInfo.KG.Service == "houndify" {
		fmt.Println("using houndify")
		houndClient := houndify.Client{
			ClientID:  uInfo.KG.APIKey,
			ClientKey: uInfo.KG.ClientKey,
		}
		houndResp := StreamAudioToHoundify(sReq, houndClient)
		spokenText, _ = ParseSpokenResponse(houndResp)
		fmt.Println(spokenText)
	} else if uInfo.KG.Service == "openai" || uInfo.KG.Service == "together" {
		var c *openai.Client
		transcribedText, err := stt.STT(sReq, false)
		if err != nil {
			spokenText = "There was an error with your L L M settings: " + err.Error()
		} else {
			var model string
			if uInfo.KG.Service == "openai" {
				fmt.Println("using openai")
				c = openai.NewClient(uInfo.KG.APIKey)
				model = openai.GPT4o
			} else {
				fmt.Println("using together")
				conf := openai.DefaultConfig(uInfo.KG.APIKey)
				conf.BaseURL = "https://api.together.xyz/v1"
				c = openai.NewClientWithConfig(conf)
				model = "meta-llama/Llama-3-70b-chat-hf"
			}
			smsg := openai.ChatCompletionMessage{
				Role: openai.ChatMessageRoleSystem,
			}
			if uInfo.KG.Prompt == "" {
				smsg.Content = defaultllmPrompt + llmPromptAddition
			} else {
				smsg.Content = uInfo.KG.Prompt + llmPromptAddition
			}
			var nChat []openai.ChatCompletionMessage
			nChat = append(nChat, smsg)
			nChat = append(nChat, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: transcribedText,
			})
			// we have 255 characters to work with. he won't speak any more than that. calculate tokens
			mTokens := int(math.Round(float64((len([]rune(smsg.Content)) / 4))) + 63)
			fmt.Println(fmt.Sprint(mTokens) + " tokens calculated")

			aireq := openai.ChatCompletionRequest{
				Model:            model,
				MaxTokens:        mTokens,
				Temperature:      1,
				TopP:             1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Messages:         nChat,
				Stream:           false,
			}
			resp, err := c.CreateChatCompletion(context.Background(), aireq)
			if err != nil {
				spokenText = "There was an error with your L L M choice: " + err.Error()
			} else {
				spokenText = removeSpecialCharacters(resp.Choices[0].Message.Content)
			}
		}
	}
	fmt.Println("Response: " + spokenText)
	kgResp := &pb.KnowledgeGraphResponse{
		QueryText:   "unknown",
		SpokenText:  spokenText,
		CommandType: "question",
		DomainsUsed: []string{"unknown"},
	}
	req.Stream.Send(kgResp)
	return &vtt.KnowledgeGraphResponse{
		Intent: kgResp,
	}, nil
}
