package common

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"strconv"
	"strings"
)

// ArgParser is a helper to read several parameters from the request
type ArgParser struct {
	violations []string
	request    *events.APIGatewayProxyRequest
}

func NewArgParser(request *events.APIGatewayProxyRequest) *ArgParser {
	return &ArgParser{
		request: request,
	}
}

func (a *ArgParser) HasViolations() bool {
	return len(a.violations) > 0
}

func (a *ArgParser) BadRequest() (Response, error) {
	return BadRequest(map[string]string{
		"error": strings.Join(a.violations, ", "),
	})
}

func (a *ArgParser) ReadPathParameterString(key string) string {
	if value, ok := a.request.PathParameters[key]; ok {
		return value
	}

	a.violations = append(a.violations, fmt.Sprintf("%s is mandatory", key))
	return ""
}

func (a *ArgParser) ReadQueryParameterInt(key string, mandatory bool) int {
	return a.readParameterInteger(a.request.QueryStringParameters, key, mandatory)
}

func (a *ArgParser) ReadPathParameterInt(key string) int {
	return a.readParameterInteger(a.request.PathParameters, key, true)
}

func (a *ArgParser) ReadQueryParameterBool(key string, mandatory bool) bool {
	return a.readParameterBool(a.request.QueryStringParameters, key, mandatory)
}

func (a *ArgParser) readParameterInteger(parameters map[string]string, key string, mandatory bool) int {
	if value, ok := parameters[key]; ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			a.violations = append(a.violations, fmt.Sprintf("%s must be a number, got '%s'", key, value))
		}

		return num
	}

	if mandatory {
		a.violations = append(a.violations, fmt.Sprintf("%s is mandatory", key))
	}

	return 0
}

func (a *ArgParser) readParameterBool(parameters map[string]string, key string, mandatory bool) bool {
	if value, ok := parameters[key]; ok {
		num, err := strconv.ParseBool(value)
		if err != nil {
			a.violations = append(a.violations, fmt.Sprintf("%s must be a number, got '%s'", key, value))
		}

		return num
	}

	if mandatory {
		a.violations = append(a.violations, fmt.Sprintf("%s is mandatory", key))
	}

	return false
}
