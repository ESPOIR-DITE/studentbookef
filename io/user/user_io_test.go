package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"studentbookef/domain"
	"testing"
	"time"
)

func TestUserLog(t *testing.T) {
	userAccount := domain.UserAccount{"espoir@admin.com", "HQfeJib6", "", time.Now()}
	result, err := UserLog(userAccount)

	fmt.Println(result)
	assert.Nil(t, err)
}
func TestReadAllLog(t *testing.T) {
	result, err := ReadAllLog()
	fmt.Println(result)
	assert.Nil(t, err)
}
func TestCreateUserPost(t *testing.T) {
	userPost := domain.UserPost{"BPF-3890174d-f67f-43fd-91f2-989d0574b476", "216093805@mycput.ac.za"}
	result, err := CreateUserPost(userPost)
	fmt.Println(result)
	assert.Nil(t, err)
}
func TestReadUserPosts(t *testing.T) {
	result, err := ReadUserPosts()
	fmt.Println(result)
	assert.Nil(t, err)
}
func TestReafAllOfUserPost(t *testing.T) {
	result, err := ReadAllOfUserPost("216093805@mycput.ac.za")
	fmt.Println(result)
	assert.Nil(t, err)
}

func TestReadUsers(t *testing.T) {
	result, err := ReadUsers()
	fmt.Println(result)
	assert.Nil(t, err)
}
