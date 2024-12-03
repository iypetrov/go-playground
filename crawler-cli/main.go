package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/disintegration/imaging"
)

func Render(w io.Writer, name string, data interface{}) error {
	return template.Must(template.ParseGlob("*.html")).ExecuteTemplate(w, name, data)
}

func main() {
	// Define the CLI command
	var cmd = &cobra.Command{
		Use:   "crawler-cli [url]",
		Short: "A web crawler that crawls web pages and extracts images",
		Args:  cobra.ExactArgs(1), // Ensure one argument is provided
		Run: func(cmd *cobra.Command, args []string) {
			url := args[0]
			outputDir := "./thumbnails"

			// Start the crawling process
			if err := CrawlAndSaveImages(url, outputDir); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Crawling and image processing completed successfully.")
		},
	}

	go func() {
		// Handle the root path and thumbnails folder
		http.Handle("/thumbnails/", http.StripPrefix("/thumbnails", http.FileServer(http.Dir("thumbnails"))))
		http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
			// Get a list of all files in the thumbnails directory
			files, err := os.ReadDir("thumbnails")
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to read thumbnails directory: %v", err), http.StatusInternalServerError)
				return
			}
		
			// Prepare a list of file names
			var fileNames []string
			for _, file := range files {
				if !file.IsDir() {
					fileNames = append(fileNames, "/thumbnails/"+file.Name())
				}
			}
		
			Render(w, "index", struct{ Files []string }{fileNames})
		})
		

		// Start the HTTP server on port 8080
		fmt.Println(http.ListenAndServe(":8080", nil))
	}()

	// Execute the CLI command
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	select {}
}

// CrawlAndSaveImages crawls a webpage for images, downloads them, and saves thumbnails.
func CrawlAndSaveImages(resourceURL, outputDir string) error {
	// Fetch the webpage content
	resp, err := http.Get(resourceURL)
	if err != nil {
		return fmt.Errorf("failed to fetch webpage: %w", err)
	}
	defer resp.Body.Close()

	// Ensure the response is HTML
	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "text/html") {
		return errors.New("the URL does not return an HTML document")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract image URLs
	imageURLs := extractImageURLs(string(body), resourceURL)

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Process each image
	for _, imageURL := range imageURLs {
		if err := downloadAndSaveThumbnail(imageURL, outputDir); err != nil {
			fmt.Printf("Failed to process image %s: %v\n", imageURL, err)
		}
	}

	return nil
}

// extractImageURLs extracts image URLs from the HTML content.
func extractImageURLs(htmlContent, baseURL string) []string {
	// Regex to match image src attributes
	imgRegex := regexp.MustCompile(`(?i)<img[^>]+src="([^"]+)"`)

	// Find all matches
	matches := imgRegex.FindAllStringSubmatch(htmlContent, -1)

	// Process matches to resolve full URLs
	var imageURLs []string
	for _, match := range matches {
		if len(match) > 1 {
			imageURL := match[1]
			// Resolve relative URLs
			if strings.HasPrefix(imageURL, "//") {
				imageURL = "http:" + imageURL
			} else if strings.HasPrefix(imageURL, "/") {
				imageURL = strings.TrimSuffix(baseURL, "/") + imageURL
			}
			imageURLs = append(imageURLs, imageURL)
		}
	}
	return imageURLs
}

// downloadAndSaveThumbnail downloads an image, creates a thumbnail, and saves it.
func downloadAndSaveThumbnail(imageURL, outputDir string) error {
	// Fetch the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Decode the image
	img, format, err := image.Decode(resp.Body)
	if err != nil {
		// Skip unsupported image formats
		return fmt.Errorf("failed to decode image (%s): %w", imageURL, err)
	}

	// Resize the image
	thumbnail := imaging.Resize(img, 200, 0, imaging.Lanczos)

	// Save the thumbnail
	filename := filepath.Base(imageURL)
	filename = strings.Split(filename, "?")[0] // Remove query parameters
	ext := filepath.Ext(filename)
	filename = strings.TrimSuffix(filename, ext) + "_thumb.jpg"

	outputPath := filepath.Join(outputDir, filename)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if format == "jpeg" || format == "jpg" {
		err = jpeg.Encode(file, thumbnail, nil)
	} else {
		err = imaging.Save(thumbnail, file.Name())
	}

	if err != nil {
		return fmt.Errorf("failed to save thumbnail: %w", err)
	}

	fmt.Printf("Thumbnail saved: %s\n", outputPath)
	return nil
}
