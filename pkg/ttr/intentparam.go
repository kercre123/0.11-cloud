package ttr

import (
	"fmt"
	"strings"

	lcztn "github.com/kercre123/0.11-cloud/pkg/localization"
	"github.com/kercre123/0.11-cloud/pkg/vars"
)

func prehistoricParamChecker(req interface{}, intent string, speechText string, botSerial string, isOpus bool) {
	var intentParam string
	var intentParamValue string
	var newIntent string
	var isParam bool
	var intentParams map[string]string
	var botLocation string = "San Francisco"
	var botUnits string = "F"
	isNew, uInfo := vars.GetUserInfo(botSerial)
	if !isNew {
		botLocation = uInfo.Location
	}
	if strings.Contains(intent, "intent_photo_take_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_ME)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_SELF)) {
			intentParam = "entity_photo_selfie"
			intentParamValue = "photo_selfie"
		} else {
			intentParam = "entity_photo_selfie"
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_eyecolor") {
		// leaving stuff like this in case someone wants to add features like this to older software
		isParam = true
		newIntent = "intent_imperative_eyecolor_specific_extend"
		intentParam = "eye_color"
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_PURPLE)) {
			intentParamValue = "COLOR_PURPLE"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_BLUE)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_SAPPHIRE)) {
			intentParamValue = "COLOR_BLUE"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_YELLOW)) {
			intentParamValue = "COLOR_YELLOW"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_TEAL)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_TEAL2)) {
			intentParamValue = "COLOR_TEAL"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_GREEN)) {
			intentParamValue = "COLOR_GREEN"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_EYE_COLOR_ORANGE)) {
			intentParamValue = "COLOR_ORANGE"
		} else {
			newIntent = intent
			intentParamValue = ""
			isParam = false
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_weather_extend") {
		isParam = true
		newIntent = intent
		condition, is_forecast, local_datetime, speakable_location_string, temperature, temperature_unit := weatherParser(speechText, botLocation, botUnits)
		intentParams = map[string]string{"condition": condition, "is_forecast": is_forecast, "local_datetime": local_datetime, "speakable_location_string": speakable_location_string, "temperature": temperature, "temperature_unit": temperature_unit}
	} else if strings.Contains(intent, "intent_imperative_volumelevel_extend") {
		isParam = true
		newIntent = intent
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_MEDIUM_LOW)) {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_2"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_LOW)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_QUIET)) {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_MEDIUM_HIGH)) {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_4"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_MEDIUM)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_NORMAL)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_REGULAR)) {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_3"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_HIGH)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_LOUD)) {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_5"
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_MUTE)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_NOTHING)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_SILENT)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_OFF)) || strings.Contains(speechText, lcztn.GetText(lcztn.STR_VOLUME_ZERO)) {
			// there is no VOLUME_0 :(
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		} else {
			intentParam = "volume_level"
			intentParamValue = "VOLUME_1"
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_names_username_extend") {
		var username string
		var nameSplitter string = ""
		isParam = true
		if !isOpus {
			newIntent = "intent_names_username"
		} else {
			newIntent = "intent_names_username_extend"
		}
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_NAME_IS)) {
			nameSplitter = lcztn.GetText(lcztn.STR_NAME_IS)
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_NAME_IS2)) {
			nameSplitter = lcztn.GetText(lcztn.STR_NAME_IS2)
		} else if strings.Contains(speechText, lcztn.GetText(lcztn.STR_NAME_IS3)) {
			nameSplitter = lcztn.GetText(lcztn.STR_NAME_IS3)
		}
		if nameSplitter != "" {
			splitPhrase := strings.SplitAfter(speechText, nameSplitter)
			username = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				username = username + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				username = username + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			fmt.Println("Name parsed from speech: " + "`" + username + "`")
			intentParam = "username"
			intentParamValue = username
			intentParams = map[string]string{intentParam: intentParamValue}
		} else {
			fmt.Println("No name parsed from speech")
			intentParam = "username"
			intentParamValue = ""
			intentParams = map[string]string{intentParam: intentParamValue}
		}
	} else if strings.Contains(intent, "intent_clock_settimer_extend") {
		isParam = true
		newIntent = intent
		timerSecs := words2num(speechText)
		fmt.Println("Seconds parsed from speech: " + timerSecs)
		intentParam = "timer_duration"
		intentParamValue = timerSecs
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_global_stop_extend") {
		isParam = true
		newIntent = "intent_global_stop"
		intentParam = "what_to_stop"
		intentParamValue = "timer"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_playmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_playmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_FOR)) {
			splitPhrase := strings.SplitAfter(speechText, lcztn.GetText(lcztn.STR_FOR))
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_message_recordmessage_extend") {
		var given_name string
		isParam = true
		newIntent = "intent_message_recordmessage"
		intentParam = "given_name"
		if strings.Contains(speechText, lcztn.GetText(lcztn.STR_FOR)) {
			splitPhrase := strings.SplitAfter(speechText, lcztn.GetText(lcztn.STR_FOR))
			given_name = strings.TrimSpace(splitPhrase[1])
			if len(splitPhrase) == 3 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2])
			} else if len(splitPhrase) == 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			} else if len(splitPhrase) > 4 {
				given_name = given_name + " " + strings.TrimSpace(splitPhrase[2]) + " " + strings.TrimSpace(splitPhrase[3])
			}
			intentParamValue = given_name
		} else {
			intentParamValue = ""
		}
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_blackjack") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "blackjack"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_fistbump") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "fist_bump"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_play_rollcube") {
		isParam = true
		newIntent = "intent_play_specific_extend"
		intentParam = "entity_behavior"
		intentParamValue = "roll_cube"
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_praise") {
		isParam = false
		newIntent = "intent_imperative_affirmative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else if strings.Contains(intent, "intent_imperative_abuse") {
		isParam = false
		newIntent = "intent_imperative_negative"
		intentParam = ""
		intentParamValue = ""
		intentParams = map[string]string{intentParam: intentParamValue}
	} else {
		newIntent = intent
		intentParam = ""
		intentParamValue = ""
		isParam = false
		intentParams = map[string]string{intentParam: intentParamValue}
	}
	IntentPass(req, newIntent, speechText, intentParams, isParam)
}
