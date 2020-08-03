package user

import (
	"errors"
	"studentbookef/api"
	"studentbookef/domain"
)

const userroleURL = api.BASE_URL + "user_role/"

func CreateUserRole(role domain.UserRole) (domain.UserRole, error) {
	entity := domain.UserRole{}
	resp, _ := api.Rest().SetBody(role).Post(userroleURL + "create")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}

func ReadUserRole(id string) (domain.UserRole, error) {
	entity := domain.UserRole{}
	resp, _ := api.Rest().Get(userroleURL + "read?id=" + id)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}

func ReadUserRoleWithEmail(email string) (domain.UserRole, error) {
	entity := domain.UserRole{}
	resp, _ := api.Rest().Get(userroleURL + "readWithEmail?id=" + email)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}

func ReadUserRoles() ([]domain.UserRole, error) {
	entity := []domain.UserRole{}
	resp, _ := api.Rest().Get(userroleURL + "reads")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func UpdateUserRole(userRole domain.UserRole) (domain.UserRole, error) {
	entity := domain.UserRole{}
	resp, _ := api.Rest().SetBody(userRole).Post(userroleURL + "update")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}

func DeleteUserRole(userRole domain.UserRole) (domain.UserRole, error) {
	entity := domain.UserRole{}
	resp, _ := api.Rest().SetBody(userRole).Post(userroleURL + "delete")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
