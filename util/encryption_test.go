package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const HASHED_VALUE = "$2a$10$aY4BQ4iw6pbXHKMs3bi.yOClCN79NP4E288GJxGo7GLJDnJON97UW"

func TestHashGenerate(t *testing.T) {
	_, err := HashGenerate("1234")
	assert.NoError(t, err)
}

func TestHashCompare(t *testing.T) {
	/* Success Case */
	fmt.Println("------- TestHashCompare Success Case -------")
	err := HashCompare(HASHED_VALUE, "1234")
	assert.NoError(t, err)

	/* Failure Case */
	fmt.Println("------- TestHashCompare Failure Case -------")
	err = HashCompare(HASHED_VALUE, "5678")
	assert.Error(t, err)

}
