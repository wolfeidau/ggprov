package ggprov

import "github.com/aws/aws-sdk-go/aws/awserr"

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
