package userusecase_test

import (
	"fmt"
	"rap-c/app/entity"

	"github.com/golang/mock/gomock"
)

type createMatcher struct {
	expected *entity.User
	want     string
}

func (r *createMatcher) Matches(x interface{}) bool {
	req, ok := x.(*entity.User)
	if !ok {
		return false
	}
	r.expected.Password = req.Password
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
	if req.PasswordMustChange != r.expected.PasswordMustChange {
		return false
	}
	if req.IsGuest != r.expected.IsGuest {
		return false
	}
	if req.CreatedBy != r.expected.CreatedBy {
		return false
	}
	if req.UpdatedBy != r.expected.UpdatedBy {
		return false
	}
	// validate password, password is auto generated, so check if not empty only
	if req.Password == "" {
		r.want = "password not empty"
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

func CreateMatcher(expected *entity.User) gomock.Matcher {
	return &createMatcher{expected: expected}
}
