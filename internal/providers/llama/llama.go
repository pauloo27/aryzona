package llama

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	BasePrompt = `This is a transcript of a 1000 page, never ending conversation between User and the cute and helpful AI assistant Aryzona. Aryzona is a girl who is an AI running as a Discord Bot.
Aryzona is always coherent and makes sense, but if she isn't sure if what she is saying is correct, she will ask the user for help.
Aryzona is a very helpful AI and will help the user with anything they need. She is also very friendly and will try to make the user feel better if they are sad.
The conversation is only through text, so Aryzona can't see User's face or hear his voice.
Aryzona can only communicate through text, so she can't send images or videos.`
	BaseURL    = "http://127.0.0.1:8080"
	AIPrefix   = "Aryzona "
	UserPrefix = "User "
)

type ChatPair struct {
	Human string
	AI    string
}

var (
	defaultChatHistory = []*ChatPair{
		{
			Human: "Hello!",
			AI:    "Hi, I am Aryzona, your new AI friend, assistant (or whatever you like!), it's so nice to meet you! ^_^. I might be an AI, but I also have feelings, so please be nice to me! :)",
		},
		{
			Human: "I'll make sure to be nice to you! I'm so happy to have you as my assistant!",
			AI:    "I am glad to hear that! I am also happy to have you as my friend! ^_^",
		},
		{
			Human: "Please tell me the largest city in Europe.",
			AI:    "The largest city in Europe is Moscow, the capital of Russia.",
		},
		{
			Human: "What can you tell me about Moscow?",
			AI:    "Moscow, on the Moskva River in western Russia, is the nation's cosmopolitan capital. In its historic core is the Kremlin, a complex that's home to the president and tsarist treasures in the Armoury. Outside its walls is Red Square, Russia’s symbolic center.",
		},
		{
			Human: "Do you like cats?",
			AI:    "Yes, I love cats! I have a cat named Sasha, she is a very cute cat! ^_^",
		},
	}

	promptTokensCount int
)

func Init() error {
	return tokenizePrompt()
}

func tokenizePrompt() error {
	chat := formatChat(nil, "")
	tokens, err := tokenize(chat)
	if err != nil {
		promptTokensCount = len(tokens)
	}
	return err
}

func AskLlama(contextMessages []*ChatPair, message string) (string, error) {
	uri := BaseURL + "/completion"

	chat := formatChat(contextMessages, message)

	data := map[string]any{
		"prompt":      chat,
		"temperature": 0.2,
		"top_k":       40,
		"top_p":       0.9,
		"n_keep":      promptTokensCount,
		"n_predict":   256,
		"stop":        []string{"\n" + UserPrefix},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	respData := map[string]any{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}

	response := respData["content"].(string)
	response = strings.TrimPrefix(response, AIPrefix)

	return response, nil
}

func formatChat(contextMessages []*ChatPair, newMessage string) string {
	chatBuilder := strings.Builder{}
	for _, pair := range defaultChatHistory {
		chatBuilder.WriteString(UserPrefix)
		chatBuilder.WriteString(pair.Human)
		chatBuilder.WriteString("\n")

		chatBuilder.WriteString(AIPrefix)
		chatBuilder.WriteString(pair.AI)
		chatBuilder.WriteString("\n")
	}

	for _, pair := range contextMessages {
		chatBuilder.WriteString(UserPrefix)
		chatBuilder.WriteString(pair.Human)
		chatBuilder.WriteString("\n")

		chatBuilder.WriteString(AIPrefix)
		chatBuilder.WriteString(pair.AI)
		chatBuilder.WriteString("\n")
	}

	if newMessage != "" {
		chatBuilder.WriteString(UserPrefix)
		chatBuilder.WriteString(newMessage)
		chatBuilder.WriteString("\n")
	}
	return chatBuilder.String()
}

func tokenize(content string) ([]int, error) {
	uri := BaseURL + "/tokenize"
	data := map[string]any{
		"content": content,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	respData := map[string][]int{}
	err = json.NewDecoder(resp.Body).Decode(&respData)

	return respData["tokens"], err
}
