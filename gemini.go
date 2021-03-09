package main

import (
	"context"
	"fmt"

	"git.sr.ht/~adnano/go-gemini"
)

func Fetch(url string) (gemini.Text, error) {
	// TODO: implement UI for TOFU instead of trusting everything
	client := gemini.Client{}

	resp, err := client.Get(context.TODO(), url)
	if resp.Status == gemini.StatusRedirect || resp.Status == gemini.StatusPermanentRedirect {
		// TODO: indicate in the UI that the redirect is happening, also avoid loops
		return Fetch(resp.Meta)
	}
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
