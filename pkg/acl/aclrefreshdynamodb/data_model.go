package aclrefreshdynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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

func marshalToken(token string, spec aclcore.RefreshTokenSpec) (map[string]types.AttributeValue, error) {
	item, err := attributevalue.MarshalMap(RefreshTokenRecord{
		TablePk:             RefreshTokenRecordPk(token),
		Email:               spec.Email.Value(),
		RefreshTokenPurpose: string(spec.RefreshTokenPurpose),
		AbsoluteExpiryTime:  spec.AbsoluteExpiryTime,
		Scopes:              strings.Join(spec.Scopes, " "),
	})
	if err == nil && len(spec.Scopes) == 0 {
		item["Scopes"] = &types.AttributeValueMemberNULL{Value: true}
	}

	return item, errors.Wrapf(err, "failed to marshal %+v", spec)
}

func unmarshalToken(item map[string]types.AttributeValue) (*aclcore.RefreshTokenSpec, error) {
	record := new(RefreshTokenRecord)
	err := attributevalue.UnmarshalMap(item, record)

	return &aclcore.RefreshTokenSpec{
		Email:               usermodel.NewUserId(record.Email),
		RefreshTokenPurpose: aclcore.RefreshTokenPurpose(record.RefreshTokenPurpose),
		AbsoluteExpiryTime:  record.AbsoluteExpiryTime,
		Scopes:              strings.Split(record.Scopes, " "),
	}, errors.Wrapf(err, "failed to unmarshal RefreshTokenRecord %+v", item)
}
