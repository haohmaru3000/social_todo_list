package common

import "fmt"

const (
	CurrentUser = "current_user"
)

type DbType int

const (
	DbTypeItem DbType = 1
	DbTypeUser DbType = 2
)

const (
	PluginDBMain = "mysql"
	PluginJWT    = "jwt"
)

func Recovery() {
	if r := recover(); r != nil {
		fmt.Println("Recovered:", r)
	}
}

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}
