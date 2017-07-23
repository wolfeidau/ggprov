package ggprov

import (
	"log"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

// IamRole IAM role
type IamRole struct {
	Arn string
}

func (r *IamRole) String() string {
	return fmt.Sprintf("{ Arn: %s }", r.Arn)
}

func newRole(roleArn *string) *IamRole {
	return &IamRole{
		Arn: aws.StringValue(roleArn),
	}
}

// CreateOrGetIamRole create or get the existing iam role by name
func (s *Svcs) CreateOrGetIamRole(roleName, description, document string) (*IamRole, error) {
	log.Println("Get service role")

	grresp, err := s.IAMAPI.GetRole(&iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil && !isNotFoundErr(err) {
		return nil, errors.Wrap(err, "Failed to get role")
	}
	// create the resource
	if isNotFoundErr(err) {
		log.Println("Creating service role for account")
		// aws iam create-role --role-name Greengrass-Service-Role \
		// --description "Allows AWS Greengrass to call AWS Services on your behalf" \
		// --assume-role-policy-document
		crresp, err := s.IAMAPI.CreateRole(&iam.CreateRoleInput{
			RoleName:                 aws.String(roleName),
			AssumeRolePolicyDocument: aws.String(document),
			Description:              aws.String(description),
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create role")
		}

		return newRole(crresp.Role.Arn), nil
	}

	return newRole(grresp.Role.Arn), nil
}
