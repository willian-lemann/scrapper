package external

import (
	"context"
	"fmt"
	"os"
	"scrapper/utils"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func OpenAi(markdown string) {
	utils.LoadEnvs()

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL("https://api.deepseek.com"))

	userPrompt := fmt.Sprintf(`Extract the following information from the markdown content below: 
	name
	address
	price
	link
	image
	area (in square meters)
	bedrooms
	type (if is a house or apartment, or other thing)
	forSale (if is for sale or rent)
	parking (if has parking or not)
	content (is the description of the property)
	photos (if has photos or not)
	agency (if has agency or not)
	bathrooms (if has bathrooms or not)
	ref (if has a reference or not)
	Markdown content:
	%s`, markdown)

	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(`You are a data extraction tool that reads markdown files. These markdown files are the result of web scraping real estate websites. Your task is to extract and return only the following information from the provided markdown content:

	name
	address
	price
	link
	image
	area (in square meters)
	bedrooms
	type (if is a house or apartment, or other thing)
	forSale (if is for sale or rent)
	parking (if has parking or not)
	content (is the description of the property)
	photos (if has photos or not)
	agency (if has agency or not)
	bathrooms (if has bathrooms or not)
	ref (if has a reference or not)
Return the extracted information in a structured JSON format.`),
			openai.UserMessage(userPrompt),
		},
		Model: "deepseek-chat",
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
