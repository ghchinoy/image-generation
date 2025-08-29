# Image Generation CLI

This command-line tool generates images using Google's generative AI models.

## Usage

```sh
go run main.go -prompt "your image prompt"
```

### Flags

*   `-prompt` (required): The text prompt for image generation.
*   `-model`: The name of the model to use. Defaults to `gemini-2.5-flash-image-preview`.
*   `-image`: Path to an image file or a GCS URI. This flag can be repeated to provide multiple images.

### Environment Variables

The following environment variables must be set:

*   `GOOGLE_CLOUD_PROJECT`: Your Google Cloud project ID.
*   `GOOGLE_CLOUD_LOCATION`: The Google Cloud location for the model.

### Example

```sh
export GOOGLE_CLOUD_PROJECT="your-gcp-project"
export GOOGLE_CLOUD_LOCATION="us-central1"

go run main.go -prompt "a cat wearing a hat" -image "cat.png"
```

### Output

The tool will save the generated images as PNG files in the current directory and print a JSON array of the file paths to the console.
