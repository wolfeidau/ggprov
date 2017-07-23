package ggprov

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/greengrass/greengrassiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iot/iotiface"
)

// Svcs holds aws session and svcs
type Svcs struct {
	GreengrassAPI greengrassiface.GreengrassAPI
	IAMAPI        iamiface.IAMAPI
	IoTAPI        iotiface.IoTAPI
	Session       *session.Session
}

// CreateSvcs create an aws session and configure the svcs
func CreateSvcs() (*Svcs, error) {
	svcs := &Svcs{}

	sess := session.Must(session.NewSession())
	svcs.GreengrassAPI = greengrass.New(sess)
	svcs.IAMAPI = iam.New(sess)
	svcs.IoTAPI = iot.New(sess)

	svcs.Session = sess

	return svcs, nil
}
