package ggprov

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/pkg/errors"
)

// Thing an aws thing
type Thing struct {
	Name string
	Arn  string
}

// IotPolicy IAM role
type IotPolicy struct {
	Name string
	Arn  string
}

// IotEndpoint IoT Service Endpoint
type IotEndpoint struct {
	Hostname string
}

// ThingCreds the certificates and keys for aws thing
type ThingCreds struct {
	CertificateArn string
	CertificatePem string
	PrivateKey     string
	PublicKey      string
}

func newThing(ctresp *iot.CreateThingOutput) *Thing {
	return &Thing{
		Name: aws.StringValue(ctresp.ThingName),
		Arn:  aws.StringValue(ctresp.ThingArn),
	}
}

func newPolicy(policyName *string, policyArn *string) *IotPolicy {
	return &IotPolicy{
		Name: aws.StringValue(policyName),
		Arn:  aws.StringValue(policyArn),
	}
}

func newEndpoint(deresp *iot.DescribeEndpointOutput) *IotEndpoint {
	return &IotEndpoint{
		Hostname: aws.StringValue(deresp.EndpointAddress),
	}
}

func newThingCreds(ckacresp *iot.CreateKeysAndCertificateOutput) *ThingCreds {
	return &ThingCreds{
		CertificateArn: aws.StringValue(ckacresp.CertificateArn),
		CertificatePem: aws.StringValue(ckacresp.CertificatePem),
		PrivateKey:     aws.StringValue(ckacresp.KeyPair.PrivateKey),
		PublicKey:      aws.StringValue(ckacresp.KeyPair.PublicKey),
	}
}

// CreateThing create an AWS IoT thing
func (s *Svcs) CreateThing(thingName string) (*Thing, error) {

	log.Println("Creating thing")

	ctresp, err := s.IoTAPI.CreateThing(&iot.CreateThingInput{
		ThingName: aws.String(thingName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create thing")
	}

	return newThing(ctresp), nil
}

// CreateKeysAndCertificates create and active thing certificates and keys
func (s *Svcs) CreateKeysAndCertificates() (*ThingCreds, error) {
	log.Println("Creating thing certificates")

	ckacresp, err := s.IoTAPI.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{
		SetAsActive: aws.Bool(true),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create thing certificates and keys")
	}

	return newThingCreds(ckacresp), nil
}

// CreateThingPolicy create a thing policy
func (s *Svcs) CreateThingPolicy(thingName string) (*IotPolicy, error) {
	log.Println("Creating thing policy for", thingName)
	return s.CreateIotPolicy(fmt.Sprintf("%s-DEPLOYMENT-IOT-Policy", thingName), ggThingPolicyDocument)
}

// CreateIotPolicy create an iam policy
func (s *Svcs) CreateIotPolicy(policyName, document string) (*IotPolicy, error) {
	log.Println("Get service policy")

	cpresp, err := s.IoTAPI.CreatePolicy(&iot.CreatePolicyInput{
		PolicyName:     aws.String(policyName),
		PolicyDocument: aws.String(document),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create policy")
	}

	return newPolicy(cpresp.PolicyName, cpresp.PolicyArn), nil
}

// AttachPrincipalPolicy attach the thing policy to the principal
func (s *Svcs) AttachPrincipalPolicy(thingCreds *ThingCreds, policy *IotPolicy) error {
	log.Println("Attach policy to principal", thingCreds.CertificateArn, policy.Name)

	_, err := s.IoTAPI.AttachPrincipalPolicy(&iot.AttachPrincipalPolicyInput{
		Principal:  aws.String(thingCreds.CertificateArn),
		PolicyName: aws.String(policy.Name),
	})

	return errors.Wrap(err, "Failed to attach Policy to Principal")
}

// AttachThingPrincipal attach thing to principal
func (s *Svcs) AttachThingPrincipal(thing *Thing, thingCreds *ThingCreds) error {
	log.Println("Attach thing to principal")

	_, err := s.IoTAPI.AttachThingPrincipal(&iot.AttachThingPrincipalInput{
		Principal: aws.String(thingCreds.CertificateArn),
		ThingName: aws.String(thing.Name),
	})

	return errors.Wrap(err, "Failed to attach Thing to Principal")
}

// GetIoTEndpoint get the iot endpoint information
func (s *Svcs) GetIoTEndpoint() (*IotEndpoint, error) {
	log.Println("Get IoT Endpoint")

	geresp, err := s.IoTAPI.DescribeEndpoint(&iot.DescribeEndpointInput{})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to retrieve endpoint")
	}

	return newEndpoint(geresp), nil
}
