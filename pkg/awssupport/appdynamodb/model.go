package appdynamodb

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// TablePk are the primary and sort keys of the table
type TablePk struct {
	PK string // PK is the Partition key ; see what's used depending on object types
	SK string // SK is the Sort key ; see what's used depending on object types
}

// MediaPrimaryKeyPK is the PK of a media, used to regroup media related information together
func MediaPrimaryKeyPK(owner string, id string) string {
	return fmt.Sprintf("%s#MEDIA#%s", owner, id)
}

// UserPk is the PK of a user, used to regroup user related information together
func UserPk(userEmail usermodel.UserId) string {
	return fmt.Sprintf("USER#%s", userEmail)
}
