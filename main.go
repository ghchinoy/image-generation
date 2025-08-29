package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"
)

type stringSlice []string

func (i *stringSlice) String() string {
	return strings.Join(*i, ", ")
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var prompt string
	flag.StringVar(&prompt, "prompt", "", "The prompt for image generation")
	var model string
	flag.StringVar(&model, "model", "gemini-2.5-flash-image-preview", "the model name, e.g. gemini-2.5-flash-image-preview")
	var images stringSlice
	flag.Var(&images, "image", "Image file path or GCS URI. Can be repeated.")
	flag.Parse()

	if prompt == "" {
		log.Fatal("The -prompt flag is required.")
	}

	ctx := context.Background()

	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	location := os.Getenv("GOOGLE_CLOUD_LOCATION")

	if project == "" || location == "" {
		log.Fatal("Please set the GOOGLE_CLOUD_PROJECT and GOOGLE_CLOUD_LOCATION environment variables.")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:  project,
		Location: location,
		Backend:  genai.BackendVertexAI,
	})
	if err != nil {
		log.Fatal(err)
	}

	config := &genai.GenerateContentConfig{}
	config.ResponseModalities = []string{"IMAGE", "TEXT"}

	var parts []*genai.Part
	parts = append(parts, genai.NewPartFromText(prompt))

	for _, imgPath := range images {
		if strings.HasPrefix(imgPath, "gs://") {
			parts = append(parts, genai.NewPartFromURI(imgPath, ""))
		} else {
			imgData, err := os.ReadFile(imgPath)
			if err != nil {
				log.Fatal(err)
			}
			parts = append(parts, genai.NewPartFromBytes(imgData, inferMimeType(imgPath)))
		}
	}

	contents := &genai.Content{Parts: parts, Role: "USER"}

	resp, err := client.Models.GenerateContent(ctx, model, []*genai.Content{contents}, config)
	if err != nil {
		log.Fatal(err)
	}

	var filePaths []string
	gentime := time.Now().Format("20060102150405")

	for _, candidate := range resp.Candidates {
		for n, part := range candidate.Content.Parts {
			log.Printf("%s", part.Text)
			if part.InlineData != nil {
				log.Printf("part %d mime-type: %s", n, part.InlineData.MIMEType)

				fileName := fmt.Sprintf("image_%s_%d.png", gentime, n)
				if err := os.WriteFile(fileName, part.InlineData.Data, 0644); err != nil {
					log.Fatal(err)
				}
				filePaths = append(filePaths, fileName)
			}
		}
	}

	jsonOutput, err := json.Marshal(filePaths)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonOutput))
}

func inferMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
