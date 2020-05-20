package book_io

import (
	"errors"
	"studentbookef/api"
	"studentbookef/domain"
)

const booklanguageURL = api.BASE_URL + "book_language/"

func CreateBookLanguage(language domain.BookLanguage) (domain.BookLanguage, error) {
	entity := domain.BookLanguage{}
	resp, _ := api.Rest().SetBody(language).Post(booklanguageURL + "create")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func DeleteBookLanguage(language domain.BookLanguage) (domain.BookLanguage, error) {
	entity := domain.BookLanguage{}
	resp, _ := api.Rest().SetBody(language).Post(booklanguageURL + "delete")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func UpdateBookLanguage(language domain.BookLanguage) (domain.BookLanguage, error) {
	entity := domain.BookLanguage{}
	resp, _ := api.Rest().SetBody(language).Post(booklanguageURL + "update")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func ReadBookLanguage(id string) (domain.BookLanguage, error) {
	entity := domain.BookLanguage{}
	resp, _ := api.Rest().Get(booklanguageURL + "read?id=" + id)
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
func ReadBookLanguages() ([]domain.BookLanguage, error) {
	entity := []domain.BookLanguage{}
	resp, _ := api.Rest().Get(booklanguageURL + "reads")
	if resp.IsError() {
		return entity, errors.New(resp.Status())
	}
	err := api.JSON.Unmarshal(resp.Body(), &entity)
	if err != nil {
		return entity, errors.New(resp.Status())
	}
	return entity, nil
}
