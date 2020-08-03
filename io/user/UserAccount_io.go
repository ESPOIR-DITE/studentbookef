package user

import (
	"errors"
	"studentbookef/api"
	"studentbookef/domain"
)

const useraccountURL = api.BASE_URL + "user_account/"

func UserLog(loginDetails domain.UserAccount) (domain.UserAccount, error) {
	var entity domain.UserAccount
	resp, _ := api.Rest().SetBody(loginDetails).Post(useraccountURL + "login")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func ReadAllLog() ([]domain.UserAccount, error) {
	var entity []domain.UserAccount
	resp, _ := api.Rest().Get(useraccountURL + "readAll")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func ReadWithpassword(code string) (domain.UserAccount, error) {
	var entity domain.UserAccount
	resp, _ := api.Rest().Get(useraccountURL + "readwithcode?id=" + code)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func UpdateUserAccount(loginDetails domain.UserAccount) (domain.UserAccount, error) {
	var entity domain.UserAccount
	resp, _ := api.Rest().SetBody(loginDetails).Post(useraccountURL + "update")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func ReadUserAccountWithEmail(email string) (domain.UserAccount, error) {
	var entity domain.UserAccount
	resp, _ := api.Rest().Get(useraccountURL + "readwithemail?id=" + email)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
