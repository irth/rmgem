package main

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
)

func Fetch(url string) (gemini.Text, error) {
	client := gemini.Client{}

	req, err := gemini.NewRequest(url)
	if err != nil {
		return nil, fmt.Errorf("in gemini.NewRequest: %w", err)
	}

	resp, err := client.Do(context.TODO(), req)
	if err != nil {
		return nil, fmt.Errorf("while making the Gemini request: %w", err)
	}
	defer resp.Body.Close()

	parsed, err := gemini.ParseText(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while parsing the Gemini response")
	}
	fmt.Printf("parsed: %+v\n", parsed)

	return parsed, nil
}
