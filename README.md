# LLM for Programmers

1. LLM as Code Assistant

2. Create LLM based product

   - Using 3rd party LLM (ChatGPT, Gemini, Claude, etc.)
   - Using local LLM

   - Examples:
     - Chat application (multimodal: text, audio and image generation)
     - Youtube video summarizer
     - Retrieval-Augmented Generation

---

## 1. LLM as Code Assistant

- Using 3rd party app (ChatGPT/Gemini etc.)

  - [ChatGPT](https://chat.openai.com/)
  - [Gemini](https://gemini.google.com/)

- Using code editor/IDE plugins (Copilot/Cody/Codium etc.)

---

### Bug fixing

```python

def add_item(item, shopping_list=[]):
    shopping_list.append(item)
    return shopping_list


list1 = add_item("apples")  # ["apples"]
list2 = add_item("bananas")  # ["bananas"]

print("List 1:", list1)
print("List 2:", list2)
```

---

### Code generation

Prompt examples:

Write a python code that tracks the CPU usage of a particular process (by PID) and graph it using matplotlib library.
Track it for 30 seconds.
And save the chart in chart.png file.

---

### Library/framework exploration

Prompt examples:

- Describe to me the top 5 most frequently used `<library>` functions and give me some examples.
- What are the core fundamental concepts of `<framework>`, explain to me as someone who never used it before.

---

### Using Code Editor/IDE plugins

- [Cody](https://sourcegraph.com/cody)
- [Copilot](https://github.com/features/copilot)

---

## 2. Create LLM Based Product

- Using 3rd party / closed weight LLM
- Using local LLM
- Usage examples:
  - Chat application (multimodal: text, audio and image generation)
  - Youtube video summarizer
  - Retrieval-Augmented Generation

---

### 2a. Using 3rd party / closed weight LLM

Closed weight LLM examples:

- [Openai (GPT 3.5, GPT 4)](https://platform.openai.com/)
- [Cohere](https://cohere.com/)
- [Claude](https://www.anthropic.com/claude)

---

#### Using OpenAI API

[See the notebook](http://localhost:8888/notebooks/1_openai/OpenAI.ipynb)

---

#### Using OpenAI API

- `model` selection
- `messages` parameter (`system`, `user`, `assistant`)
- `temperature`
- `max_tokens`
- `stream`

---

### 2b. Using Local LLM

---

### 2b. Using Local LLM

- LLM is just files.
- LLM is just a git folder.

See the available open-weight LLMs [here](https://huggingface.co/models?pipeline_tag=text-generation&sort=trending)

---

#### Anatomy of LLM

- Base Model (foundation model), e.g: `Llama-3`, `Phi-3`, `OpenELM`, `Mistral`
- Parameter size, e.g: `7B`, `8B`, `70B`
- Context size, e.g: `262k`, `1048k`
- Fine-tuned data: e.g: `instruct`, `chat`, `chinese-chat`

See leaderboard [here](https://huggingface.co/spaces/lmsys/chatbot-arena-leaderboard)

---

#### Deploying LLM using [Ollama](https://ollama.com/)

1. Install ollama
2. Run server: `ollama serve`
3. Pull model: `ollama pull llama3`

[See the notebook](http://localhost:8888/notebooks/2_local-llm/Local%20LLM.ipynb)

---

### 2c. LLM Usage Example - Chat Application with Audio and Image Generation

- OpenAI API => [Text/Image generation](http://localhost:8888/notebooks/4_openai_image-gen/image.ipynb))
- Ollama (Text)
- [Prosa TTS => [Audio generation]](http://localhost:8888/notebooks/3_prosa-speech-api/speech.ipynb)

---

Run the app:

```bash
cd app
go run ./cmd/*
```

Open [this address](http://localhost:3000)

---

### 2c. LLM Usage Example - Youtube Summarizer

We are using:

- OpenAI API / Ollama (Summary generation)
- whisper.cpp (Audio transcription, audio -> text)
- yt-dlp (Youtube video downloader)

[See the notebook](http://localhost:8888/notebooks/5_whispercpp-yt_dlp%2Fwhisper_cpp-ytdlp.ipynb)

---

then open [this address](http://localhost:3000/summary)

---

### 2c. RAG

[See the notebook](http://localhost:8888/notebooks/6_anything_llm%2Fanything_llm.ipynb)
