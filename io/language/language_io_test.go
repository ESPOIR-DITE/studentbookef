package language

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"studentbookef/domain"
	"testing"
)

func TestCreateLanguage(t *testing.T) {
	language := domain.Language{"", "Swahili"}
	result, err := CreateLanguage(language)
	assert.Nil(t, err)
	fmt.Println(result)
}
func TestReadLanguage(t *testing.T) {
	result, err := ReadLanguage("")
	assert.Nil(t, err)
	fmt.Println(result)
}
func TestDeleteLanguage(t *testing.T) {
	language := domain.Language{"LF-83cbc07a-6da1-494a-9331-815c7c4a7a51", "English"}
	result, err := DeleteLanguage(language)
	assert.Nil(t, err)
	fmt.Println(result)
}
func TestReadLanguages(t *testing.T) {
	result, err := ReadLanguages()
	assert.Nil(t, err)
	fmt.Println(result)
}
