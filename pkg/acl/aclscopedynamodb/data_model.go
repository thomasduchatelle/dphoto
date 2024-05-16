package aclscopedynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"strings"
	"time"
)

const (
	scopePrefix = "SCOPE#"
)

type ScopeRecord struct {
	appdynamodb.TablePk
	Type          string
	GrantedAt     time.Time
	GrantedTo     string
	ResourceOwner string
	ResourceId    string
	ResourceName  string
}

func ScopeRecordPk(user usermodel.UserId, scopeType string, owner ownermodel.Owner, id string) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.UserPk(user),
		SK: fmt.Sprintf("%s%s#%s#%s", scopePrefix, scopeType, owner, id),
	}
}

func MarshalScopeId(id aclcore.ScopeId) map[string]types.AttributeValue {
	pk := ScopeRecordPk(id.GrantedTo, string(id.Type), id.ResourceOwner, id.ResourceId)
	return map[string]types.AttributeValue{
		"PK": dynamoutils.AttributeValueMemberS(pk.PK),
		"SK": dynamoutils.AttributeValueMemberS(pk.SK),
	}
}

func MarshalScope(scope aclcore.Scope) (map[string]types.AttributeValue, error) {
	if err := scope.GrantedTo.IsValid(); err != nil {
		return nil, errors.Wrapf(err, "GrantedTo is mandatory to store a scope")
	}
	if isBlank(string(scope.Type)) {
		return nil, errors.New("Type is mandatory to store a scope")
	}
	if err := scope.ResourceOwner.IsValid(); err != nil {
		return nil, errors.Wrapf(err, "ResourceOwner is mandatory")
	}

	return attributevalue.MarshalMap(&ScopeRecord{
		TablePk:       ScopeRecordPk(scope.GrantedTo, string(scope.Type), scope.ResourceOwner, scope.ResourceId),
		Type:          string(scope.Type),
		GrantedAt:     scope.GrantedAt,
		GrantedTo:     scope.GrantedTo.Value(),
		ResourceOwner: scope.ResourceOwner.Value(),
		ResourceId:    scope.ResourceId,
		ResourceName:  scope.ResourceName,
	})
}

func UnmarshalScope(attributes map[string]types.AttributeValue) (*aclcore.Scope, error) {
	record := new(ScopeRecord)
	err := attributevalue.UnmarshalMap(attributes, record)

	return &aclcore.Scope{
			Type:          aclcore.ScopeType(record.Type),
			GrantedAt:     record.GrantedAt,
			GrantedTo:     usermodel.UserId(record.GrantedTo),
			ResourceOwner: ownermodel.Owner(record.ResourceOwner),
			ResourceId:    record.ResourceId,
			ResourceName:  record.ResourceName,
		},
		errors.Wrapf(err, "failed to unmarshal %+v", attributes)
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return value == "" || strings.Trim(value, " ") == ""
}
