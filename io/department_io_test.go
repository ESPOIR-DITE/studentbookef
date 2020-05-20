package io

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"studentbookef/domain"
	"testing"
)

func TestCreateDepartment(t *testing.T) {
	entity := domain.Department{"", "History", "technic"}
	result, err := CreateDepartment(entity)
	assert.Nil(t, err)
	fmt.Println("result :", result)
}
func TestReadDepartment(t *testing.T) {
	result, err := ReadDepartment("DF-f58e0739-cbf4-4f9f-aa02-a69862c71dae")
	assert.Nil(t, err)
	fmt.Println("result :", result)
}
func TestReadDepartments(t *testing.T) {
	result, err := ReadDepartments()
	assert.Nil(t, err)
	fmt.Println("result :", result)
}
func TestDeleteDepartment(t *testing.T) {
	entity := domain.Department{"DF-8d4a36ab-92ae-49ed-be78-c990ba5507e7", "IT", "technic"}
	result, err := DeleteDepartment(entity)

	assert.Nil(t, err)
	fmt.Println("result :", result)
}
func TestUpdateDepartment(t *testing.T) {
	entity := domain.Department{"", "IT", "technic"}
	result, err := DeleteDepartment(entity)

	assert.Nil(t, err)
	fmt.Println("result :", result)
}
