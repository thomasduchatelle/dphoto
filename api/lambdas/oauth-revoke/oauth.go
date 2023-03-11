package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
)

var (
	logout LogoutAdapter
)

type LogoutAdapter interface {
	RevokeSession(refreshToken string) error
}

type RevokeSessionDTO struct {
	RefreshToken string `json:"refreshToken"`
}

func Handler(request RevokeSessionDTO) (common.Response, error) {
	err := logout.RevokeSession(request.RefreshToken)
	if err != nil {
		return common.InternalError(err)
	}

	return common.Response{
		StatusCode: 204,
	}, nil
}

func main() {
	logout = common.NewLogout()

	lambda.Start(Handler)
}
