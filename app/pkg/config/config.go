package config

type Configuration struct {
	ChatModel     string  `form:"chat_model"`
	OpenaiAPIKey  *string `form:"openai_api_key"`
	OllamaModel   *string `form:"ollama_model"`
	ImageModel    string  `form:"image_model"`
	UseAudio      bool    `form:"use_audio"`
	AudioModel    string  `form:"audi_model"`
	ProsaTTSModel string  `form:"prosa_tts_model"`
}

func DefaultConfiguration() *Configuration {
	defaultOpenAIKey := ""
	defaultOllamaModel := "llama3"
	defaultImageModel := "dall-e-3"

	defaultProsaTTSMOdel := "tts-ghifari-professional"

	return &Configuration{
		ChatModel:     "ollama",
		OpenaiAPIKey:  &defaultOpenAIKey,
		OllamaModel:   &defaultOllamaModel,
		ImageModel:    defaultImageModel,
		UseAudio:      false,
		AudioModel:    "prosa",
		ProsaTTSModel: defaultProsaTTSMOdel,
	}
}
