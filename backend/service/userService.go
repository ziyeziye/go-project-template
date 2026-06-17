package service

import "go-project-template/app/errno"

var UserSvr = UserService{}

type UserService struct{}

func (s *UserService) TestAbortError() error {
	// return nil
	return errno.ServerError
}
