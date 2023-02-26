package aclrefreshdynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"strings"
	"time"
)

const (
	refreshTokenPrefix = "REFRESH#"
	skValue            = "#REFRESH_SPEC"
)

type RefreshTokenRecord struct {
	appdynamodb.TablePk
	Email               string
	RefreshTokenPurpose string
	AbsoluteExpiryTime  time.Time
	Scopes              string // Scopes is space separated
}

func RefreshTokenRecordPk(refreshToken string) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: refreshTokenPrefix + refreshToken,
		SK: skValue,
	}
}

func marshalToken(token string, spec aclcore.RefreshTokenSpec) (map[string]*dynamodb.AttributeValue, error) {
	item, err := dynamodbattribute.MarshalMap(RefreshTokenRecord{
		TablePk:             RefreshTokenRecordPk(token),
		Email:               spec.Email,
		RefreshTokenPurpose: string(spec.RefreshTokenPurpose),
		AbsoluteExpiryTime:  spec.AbsoluteExpiryTime,
		Scopes:              strings.Join(spec.Scopes, " "),
	})
	return item, errors.Wrapf(err, "failed to marshal %+v", spec)
}

func unmarshalToken(item map[string]*dynamodb.AttributeValue) (*aclcore.RefreshTokenSpec, error) {
	record := new(RefreshTokenRecord)
	err := dynamodbattribute.UnmarshalMap(item, record)

	return &aclcore.RefreshTokenSpec{
		Email:               record.Email,
		RefreshTokenPurpose: aclcore.RefreshTokenPurpose(record.RefreshTokenPurpose),
		AbsoluteExpiryTime:  record.AbsoluteExpiryTime,
		Scopes:              strings.Split(record.Scopes, " "),
	}, errors.Wrapf(err, "failed to unmarshal RefreshTokenRecord %+v", item)
}
