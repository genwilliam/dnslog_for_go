package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	t.Run("first page", func(t *testing.T) {
		result := Paginate(data, 1, 3)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("middle page", func(t *testing.T) {
		result := Paginate(data, 2, 4)
		assert.Equal(t, []int{5, 6, 7, 8}, result)
	})

	t.Run("last page partial", func(t *testing.T) {
		result := Paginate(data, 4, 3)
		assert.Equal(t, []int{10}, result)
	})

	t.Run("out of range", func(t *testing.T) {
		result := Paginate(data, 5, 5)
		assert.Equal(t, []int{}, result)
	})
}
