// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package user

import (
	"main/app/gateway/api/user"
	v1 "main/app/user-service/api/user/v1"
	"main/utility/injection"
)

type ControllerV1 struct {
	userSvc v1.UserClient
}

func NewV1() user.IUserV1 {
	return &ControllerV1{
		userSvc: injection.MustInvoke[v1.UserClient](),
	}
}
