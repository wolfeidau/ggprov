package ggprov

import (
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func isNotFoundErr(err error) bool {
	return isAwsReqErrStatusCode(err, 404)
}

func isAwsReqErrStatusCode(err error, statusCode int) bool {
	if _, ok := err.(awserr.Error); ok {
		if reqErr, ok := err.(awserr.RequestFailure); ok {
			// A service error occurred
			return reqErr.StatusCode() == statusCode
		}
	}
	return false
}

// DoClose close the "closer" and log any errors
func DoClose(output io.Closer) {
	err := output.Close()
	if err != nil {
		log.Println("Failed to close http output", err)
	}
}
