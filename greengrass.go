package ggprov

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/pkg/errors"
)

const ggRoleDocument = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "greengrass.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

const ggThingPolicyDocument = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "iot:*",
        "greengrass:*"
      ],
      "Resource": "*"
    }
  ]
}`

const (
	ggRoleName        = "Greengrass-Service-Role"
	ggRoleDescription = "Allows AWS Greengrass to call AWS Services on your behalf"
	ggPolicyArn       = "arn:aws:iam::aws:policy/service-role/AWSGreengrassResourceAccessRolePolicy"
)

// CreateOrGetServiceRoleForAccount create or get the service role
func (s *Svcs) CreateOrGetServiceRoleForAccount() (*IamRole, error) {

	var role *IamRole

	resp, err := s.GreengrassAPI.GetServiceRoleForAccount(&greengrass.GetServiceRoleForAccountInput{})
	if err != nil && !isNotFoundErr(err) {
		return nil, errors.Wrap(err, "Failed to get service role for account")
	}
	// create the resource
	if isNotFoundErr(err) {
		role, err = s.CreateOrGetIamRole(ggRoleName, ggRoleDescription, ggRoleDocument)
		if err != nil {
			return nil, err
		}

		// aws iam attach-role-policy --policy-arn arn:aws:iam::aws:policy/servicerole/AWSGreengrassResourceAccessRolePolicy \
		// --role-name $IAMROLENAME
		log.Println("Attach service role for account")

		_, err = s.IAMAPI.AttachRolePolicy(&iam.AttachRolePolicyInput{
			PolicyArn: aws.String(ggPolicyArn),
			RoleName:  aws.String(ggRoleName),
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to attach role policy")
		}

		log.Println("Associated Service Role to Account")
		_, err = s.GreengrassAPI.AssociateServiceRoleToAccount(&greengrass.AssociateServiceRoleToAccountInput{
			RoleArn: aws.String(role.Arn),
		})
		if err != nil {
			return nil, errors.Wrap(err, "Failed to associate service role to account")
		}

	} else {
		role = newRole(resp.RoleArn)
	}

	return role, nil
}
