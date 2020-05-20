package language

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"studentbookef/domain"
	"testing"
)

func TestCreateLanguage(t *testing.T) {
	language := domain.Language{"", "Zulu"}
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
	language := domain.Language{"LF-914cad3d-cbf2-423e-abfe-f79b3c42cea9", "English"}
	result, err := DeleteLanguage(language)
	assert.Nil(t, err)
	fmt.Println(result)
}
func TestReadLanguages(t *testing.T) {
	result, err := ReadLanguages()
	assert.Nil(t, err)
	fmt.Println(result)
}
