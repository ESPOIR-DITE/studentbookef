package misc

import (
	"fmt"
	"studentbookef/io/user"
)

func CheckAdmin(email string) bool {
	if email == "" {
		fmt.Println(" THe context is empty")
		return false
	}
	userAccount, err := user.ReadUserAccountWithEmail(email)
	if err != nil {
		fmt.Println(err, " could not read UserAccount")
		return false
	}
	userRole, err := user.ReadUserRole(userAccount.RoleId)
	if err != nil {
		fmt.Println(err, " could not read UserRole")
		return false
	}
	if userRole.Role == "admin" || userRole.Role == "controller" {
		return true
	}
	return false
}
