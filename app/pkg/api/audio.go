package api

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const TTSApiURL = "https://api.prosa.ai/v2/speech/tts"

var TTSApiKey string = os.Getenv("TTS_API_KEY")

type ProsaTTSReqConfig struct {
	Model       string `json:"model"`
	Wait        bool   `json:"wait"`
	AudioFormat string `json:"audio_format"`
}

type Text struct {
	Text string `json:"text"`
}

type ProsaTTSReq struct {
	Config  ProsaTTSReqConfig `json:"config"`
	Request Text              `json:"request"`
}

type ProsaTTSResponseResultData struct {
	Data string `json:"data"`
}

type ProsaTTSResponse struct {
	Status string                     `json:"string"`
	Result ProsaTTSResponseResultData `json:"result"`
}

func ProsaAudio(text string, model string) (string, error) {
	filename := fmt.Sprintf("audios/audio-%x.webm", md5.Sum([]byte(model+text)))

	// if file already exist, dont call the API
	_, err := os.Stat(filename)
	if err == nil {
		return filename, nil
	}

	payload := ProsaTTSReq{
		Config: ProsaTTSReqConfig{
			Model:       model,
			Wait:        true,
			AudioFormat: "opus",
		},
		Request: Text{
			Text: text,
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", TTSApiURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", TTSApiKey)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", errors.New("error doing the request")
	}

	dec := json.NewDecoder(res.Body)

	var response ProsaTTSResponse
	if err = dec.Decode(&response); err != nil {
		return "", errors.New("error decoding the json")
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Result.Data)
	if err != nil {
		return "", errors.New("error decoding the base64")
	}
	err = os.WriteFile(filename, decoded, 0644)
	if err != nil {
		return "", errors.New("error writing file")
	}

	return filename, nil
}
