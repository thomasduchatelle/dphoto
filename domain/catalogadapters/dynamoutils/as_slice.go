package dynamoutils

import "github.com/aws/aws-sdk-go/service/dynamodb"

// AsSlice read the stream fully and returns a slice.
func AsSlice(stream Stream) (values []map[string]*dynamodb.AttributeValue, err error) {
	for stream.HasNext() {
		values = append(values, stream.Next())
	}
	err = stream.Error()

	return
}
