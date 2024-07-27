package userusecase_test

import (
	"fmt"
	databaseentity "rap-c/app/entity/database-entity"

	"github.com/golang/mock/gomock"
)

type createMatcher struct {
	expected *databaseentity.User
	want     string
}

func (r *createMatcher) Matches(x interface{}) bool {
	req, ok := x.(*databaseentity.User)
	if !ok {
		return false
	}
	r.expected.Password = req.Password
	r.expected.Token = req.Token
	// validate each field except Password
	if req.Username != r.expected.Username {
		return false
	}
	if req.FullName != r.expected.FullName {
		return false
	}
	if req.Email != r.expected.Email {
		return false
	}
	if req.Disabled != r.expected.Disabled {
		return false
	}
	if req.PasswordMustChange != r.expected.PasswordMustChange {
		return false
	}
	if req.IsGuest != r.expected.IsGuest {
		return false
	}
	if req.CreatedByDB != r.expected.CreatedByDB {
		return false
	}
	if req.UpdatedByDB != r.expected.UpdatedByDB {
		return false
	}
	// validate password, password is auto generated, so check if not empty only
	if req.Password == "" {
		r.want = "password not empty"
		return false
	}
	// validate token, token is auto generated, so check if not empty only
	if req.Token == "" {
		r.want = "token not empty"
		return false
	}
	return true
}

func (r *createMatcher) String() string {
	if r.want != "" {
		return r.want
	}
	return fmt.Sprintf("%v (%T)", r.expected, r.expected)
}

func CreateMatcher(expected *databaseentity.User) gomock.Matcher {
	return &createMatcher{expected: expected}
}
