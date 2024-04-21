package dynamoutilsv2

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

var (
	EmptyString = &types.AttributeValueMemberS{Value: ""}
)

func AttributeValueMemberS(value string) types.AttributeValue {
	return &types.AttributeValueMemberS{Value: value}
}
