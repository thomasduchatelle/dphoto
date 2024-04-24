package dynamoutils

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AsSlice read the stream fully and returns a slice.
func AsSlice(stream Stream) (values []map[string]types.AttributeValue, err error) {
	for stream.HasNext() {
		values = append(values, stream.Next())
	}
	err = stream.Error()

	return
}
