package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Response events.APIGatewayProxyResponse

// NewJsonResponse serialises body into JSON and create a Response containing it as body.
func NewJsonResponse(code int, body interface{}, headers map[string]string) (Response, error) {
	bodyInJson, err := json.Marshal(body)
	if err != nil {
		err = errors.Wrapf(err, "failed to serialise in JSON body %+v", body)
		log.WithError(err).Errorf("serialisation failed")
		return Response{
			StatusCode: 500,
			Body:       fmt.Sprintf("serialisation failed: %s", err.Error()),
		}, nil
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, bodyInJson)

	var allHeaders = map[string]string{
		"Content-Type": "application/json",
	}
	for k, v := range headers {
		allHeaders[k] = v
	}

	return Response{
		StatusCode: code,
		Body:       buf.String(),
		Headers:    allHeaders,
	}, nil
}

// InternalError logs the error and create a 500 error response
func InternalError(cause error) (Response, error) {
	log.WithError(cause).Errorf("internal error")

	return NewJsonResponse(500, map[string]interface{}{
		"error": cause.Error(),
	}, nil)
}

func Ok(body interface{}) (Response, error) {
	return NewJsonResponse(200, body, nil)
}

func BadRequest(body interface{}) (Response, error) {
	return NewJsonResponse(400, body, nil)
}

func NotFound(body interface{}) (Response, error) {
	return NewJsonResponse(404, body, nil)
}
