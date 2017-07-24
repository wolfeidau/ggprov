package ggprov

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// DownloadFromURL Download file from URL to location
func DownloadFromURL(url, targetPath string) error {

	log.Println("Downloading", url, "to", targetPath)

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	log.Println("Downloading File", fileName)

	output, err := os.Create(targetPath)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}

	defer DoClose(output)

	response, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "Failed to download file")
	}

	defer DoClose(response.Body)

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to save file")
	}

	log.Println("Wrote", n, "to", targetPath)

	return nil
}
