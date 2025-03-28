package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PostPrompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type PostShow struct {
	Model string `json:"model"`
}

// const LargeLanguageModel = "deepseek-r1:14b"
const LargeLanguageModel = "gemma3:12b"

func main() {
	postShow := PostShow{
		Model: LargeLanguageModel,
	}
	// Create the data to send in the POST request
	postData := PostPrompt{
		Model:  LargeLanguageModel,
		Prompt: "why is the sky blue?",
		Stream: false,
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(postShow)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:11434/api/show", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	// print entire response
	totalResponse := ""
	for {
		var chunk map[string]any
		if err := decoder.Decode(&chunk); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		if modelInfo, ok := chunk["model_info"]; ok {
			// Parse model_info as JSON and pretty-print it
			modelInfoBytes, err := json.MarshalIndent(modelInfo, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			totalResponse += string(modelInfoBytes) + "\n"
		}
	}

	fmt.Println(totalResponse)

	// Convert the data to JSON
	jsonData, err = json.Marshal(postData)
	if err != nil {
		log.Fatal(err)
	}
	// Send a POST request
	resp, err = http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// Stream and collect the "response" field from the response body

	totalResponse = ""

	fmt.Println("Response Status:", resp.Status)
	decoder = json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}
		if response, ok := chunk["response"]; ok {
			totalResponse += response.(string)
		}
	}
	fmt.Println(totalResponse)
}
