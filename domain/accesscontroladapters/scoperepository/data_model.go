package scoperepository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"strings"
	"time"
)

const (
	scopePrefix = "SCOPE#"
)

type ScopeRecord struct {
	catalogdynamo.TablePk
	Type          string
	GrantedAt     time.Time
	GrantedTo     string
	ResourceOwner string
	ResourceId    string
	ResourceName  string
}

func ScopeRecordPk(user, scopeType, owner, id string) catalogdynamo.TablePk {
	return catalogdynamo.TablePk{
		PK: userPk(user),
		SK: fmt.Sprintf("%s%s#%s#%s", scopePrefix, scopeType, owner, id),
	}
}

func userPk(user string) string {
	return fmt.Sprintf("USER#%s", user)
}

func MarshalScopeId(id accesscontrol.ScopeId) map[string]*dynamodb.AttributeValue {
	pk := ScopeRecordPk(id.GrantedTo, string(id.Type), id.ResourceOwner, id.ResourceId)
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: &pk.PK},
		"SK": {S: &pk.SK},
	}
}

func MarshalScope(scope accesscontrol.Scope) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(scope.GrantedTo) {
		return nil, errors.New("GrantedTo is mandatory to store a scope")
	}
	if isBlank(string(scope.Type)) {
		return nil, errors.New("Type is mandatory to store a scope")
	}
	if isBlank(scope.ResourceId) {
		return nil, errors.New("ResourceId is mandatory")
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

func UnmarshalScopeId(attributes map[string]*dynamodb.AttributeValue) (*accesscontrol.ScopeId, error) {
	record := new(ScopeRecord)
	err := dynamodbattribute.UnmarshalMap(attributes, record)

	return &accesscontrol.ScopeId{
			Type:          accesscontrol.ScopeType(record.Type),
			GrantedTo:     record.GrantedTo,
			ResourceOwner: record.ResourceOwner,
			ResourceId:    record.ResourceId,
		},
		errors.Wrapf(err, "failed to unmarshal %+v", attributes)
}

func UnmarshalScope(attributes map[string]*dynamodb.AttributeValue) (*accesscontrol.Scope, error) {
	record := new(ScopeRecord)
	err := dynamodbattribute.UnmarshalMap(attributes, record)

	return &accesscontrol.Scope{
			Type:          accesscontrol.ScopeType(record.Type),
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
