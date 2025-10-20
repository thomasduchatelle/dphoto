package aclcognitoadapter

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type CognitoRepositoryAdapter struct {
	client     *cognitoidentityprovider.Client
	userPoolId string
}

func New(cfg aws.Config, userPoolId string) (aclcore.CognitoRepository, error) {
	if userPoolId == "" {
		return nil, errors.New("userPoolId is required")
	}

	return &CognitoRepositoryAdapter{
		client:     cognitoidentityprovider.NewFromConfig(cfg),
		userPoolId: userPoolId,
	}, nil
}

func Must(repository aclcore.CognitoRepository, err error) aclcore.CognitoRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

func (c *CognitoRepositoryAdapter) CreateUser(ctx context.Context, email usermodel.UserId, group aclcore.CognitoUserGroup) error {
	if err := email.IsValid(); err != nil {
		return errors.Wrapf(err, "invalid email address: %s", email)
	}

	// Check if user already exists
	exists, err := c.UserExists(ctx, email)
	if err != nil {
		return err
	}

	if !exists {
		// Create the user in Cognito
		_, err = c.client.AdminCreateUser(ctx, &cognitoidentityprovider.AdminCreateUserInput{
			UserPoolId: aws.String(c.userPoolId),
			Username:   aws.String(email.Value()),
			UserAttributes: []types.AttributeType{
				{
					Name:  aws.String("email"),
					Value: aws.String(email.Value()),
				},
				{
					Name:  aws.String("email_verified"),
					Value: aws.String("true"),
				},
			},
			MessageAction: types.MessageActionTypeSuppress, // Don't send welcome email as user authenticates via Google SSO
		})
		if err != nil {
			return errors.Wrapf(err, "failed to create user %s in Cognito", email)
		}

		log.WithField("Email", email).Infof("Created new Cognito user: %s", email)
	}

	// Add user to the specified group
	_, err = c.client.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(email.Value()),
		GroupName:  aws.String(string(group)),
	})
	if err != nil {
		// Check if user is already in the group
		var alreadyInGroup *types.ResourceNotFoundException
		if errors.As(err, &alreadyInGroup) {
			log.WithFields(log.Fields{
				"Email": email,
				"Group": group,
			}).Debug("User already in group")
			return nil
		}
		return errors.Wrapf(err, "failed to add user %s to group %s", email, group)
	}

	log.WithFields(log.Fields{
		"Email": email,
		"Group": group,
	}).Infof("Added user to Cognito group: %s -> %s", email, group)

	return nil
}

func (c *CognitoRepositoryAdapter) UserExists(ctx context.Context, email usermodel.UserId) (bool, error) {
	if err := email.IsValid(); err != nil {
		return false, errors.Wrapf(err, "invalid email address: %s", email)
	}

	_, err := c.client.AdminGetUser(ctx, &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(c.userPoolId),
		Username:   aws.String(email.Value()),
	})
	if err != nil {
		var userNotFound *types.UserNotFoundException
		if errors.As(err, &userNotFound) {
			return false, nil
		}
		return false, errors.Wrapf(err, "failed to check if user %s exists", email)
	}

	return true, nil
}
