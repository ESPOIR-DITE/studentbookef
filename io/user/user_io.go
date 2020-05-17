package user

import (
	"errors"
	"studentbookef/api"
	"studentbookef/domain"
)

const userURL = api.BASE_URL + "user/"

func CreateUser(user domain.User) (domain.User, error) {
	var entity domain.User
	resp, _ := api.Rest().SetBody(user).Post(userURL + "create")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}

func ReadUser(email string) (domain.User, error) {
	var entity domain.User
	resp, _ := api.Rest().Get(userURL + "read?id=" + email)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func UpdateUser(user domain.User) (domain.User, error) {
	var entity domain.User
	resp, _ := api.Rest().SetBody(user).Post(userURL + "update")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
