package dynamotestutils

import "github.com/stretchr/testify/assert"

func (d *DynamodbTestContext) Must(err error) {
	if !assert.NoError(d.T, err, "Must got an error") {
		assert.FailNow(d.T, err.Error())
	}
}

func (d *DynamodbTestContext) MustBool(value bool, err error) bool {
	if !assert.NoError(d.T, err, "Must got an error") {
		assert.FailNow(d.T, err.Error())
	}

	return value
}
