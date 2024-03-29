package aclscopedynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
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

func ScopeRecordPk(user, scopeType, owner, id string) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.UserPk(user),
		SK: fmt.Sprintf("%s%s#%s#%s", scopePrefix, scopeType, owner, id),
	}
}

func MarshalScopeId(id aclcore.ScopeId) map[string]*dynamodb.AttributeValue {
	pk := ScopeRecordPk(id.GrantedTo, string(id.Type), id.ResourceOwner, id.ResourceId)
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: &pk.PK},
		"SK": {S: &pk.SK},
	}
}

func MarshalScope(scope aclcore.Scope) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(scope.GrantedTo) {
		return nil, errors.New("GrantedTo is mandatory to store a scope")
	}
	if isBlank(string(scope.Type)) {
		return nil, errors.New("Type is mandatory to store a scope")
	}
	if isBlank(scope.ResourceOwner) {
		return nil, errors.New("ResourceOwner is mandatory")
	}

	return dynamodbattribute.MarshalMap(&ScopeRecord{
		TablePk:       ScopeRecordPk(scope.GrantedTo, string(scope.Type), scope.ResourceOwner, scope.ResourceId),
		Type:          string(scope.Type),
		GrantedAt:     scope.GrantedAt,
		GrantedTo:     scope.GrantedTo,
		ResourceOwner: scope.ResourceOwner,
		ResourceId:    scope.ResourceId,
		ResourceName:  scope.ResourceName,
	})
}

func UnmarshalScope(attributes map[string]*dynamodb.AttributeValue) (*aclcore.Scope, error) {
	record := new(ScopeRecord)
	err := dynamodbattribute.UnmarshalMap(attributes, record)

	return &aclcore.Scope{
			Type:          aclcore.ScopeType(record.Type),
			GrantedAt:     record.GrantedAt,
			GrantedTo:     record.GrantedTo,
			ResourceOwner: record.ResourceOwner,
			ResourceId:    record.ResourceId,
			ResourceName:  record.ResourceName,
		},
		errors.Wrapf(err, "failed to unmarshal %+v", attributes)
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return value == "" || strings.Trim(value, " ") == ""
}
