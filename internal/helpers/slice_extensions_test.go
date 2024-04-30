package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceContains(t *testing.T) {
	type TestingMockStruct struct {
		Name string
	}

	type TestingMockString string

	t.Run("given string slice should find object", func(t *testing.T) {
		t.Parallel()
		asset := []string{"a", "b", "c"}

		act := Contains(asset, "a")

		assert.True(t, act)
	})
	t.Run("given bool slice should find object", func(t *testing.T) {
		t.Parallel()
		asset := []bool{true, true, true}
		act := Contains(asset, true)

		assert.True(t, act)
	})
	t.Run("given int slice should find object", func(t *testing.T) {
		t.Parallel()
		asset := []int{1, 2, 3}
		act := Contains(asset, 3)

		assert.True(t, act)
	})
	t.Run("given float slice should find object", func(t *testing.T) {
		t.Parallel()
		asset := []float64{1.1, 2.1, 3.1}
		act := Contains(asset, 3.1)

		assert.True(t, act)
	})
	t.Run("given custom type slice should find object", func(t *testing.T) {
		t.Parallel()

		asset := []*TestingMockStruct{&TestingMockStruct{Name: "Pawel"}, &TestingMockStruct{Name: "Gawel"}}
		act := Contains(asset, asset[1])

		assert.True(t, act)
	})
	t.Run("given custom value type slice should find object", func(t *testing.T) {
		t.Parallel()

		asset := []TestingMockStruct{TestingMockStruct{Name: "Pawel"}, TestingMockStruct{Name: "Gawel"}}
		act := Contains(asset, TestingMockStruct{Name: "Pawel"})

		assert.True(t, act)
	})
	t.Run("given custom value string slice should find object", func(t *testing.T) {
		t.Parallel()
		var Z TestingMockString = "Z"
		var W TestingMockString = "W"
		var D TestingMockString = "D"
		asset := []TestingMockString{Z, W, D}
		act := Contains(asset, D)

		assert.True(t, act)
	})
	t.Run("given string slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []string{"a", "b", "c"}

		act := Contains(asset, "d")

		assert.False(t, act)
	})
	t.Run("given bool slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []bool{true, true, true}
		act := Contains(asset, false)

		assert.False(t, act)
	})
	t.Run("given int slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []int{1, 2, 3}
		act := Contains(asset, 4)

		assert.False(t, act)
	})
	t.Run("given float slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []float64{1.1, 2.1, 3.1}
		act := Contains(asset, 3)

		assert.False(t, act)
	})
	t.Run("given custom type slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []*TestingMockStruct{&TestingMockStruct{Name: "Pawel"}, &TestingMockStruct{Name: "Gawel"}}
		act := Contains(asset, &TestingMockStruct{Name: "Jacek"})

		assert.False(t, act)
	})
	t.Run("given custom value type slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		asset := []TestingMockStruct{TestingMockStruct{Name: "Pawel"}, TestingMockStruct{Name: "Gawel"}}
		act := Contains(asset, TestingMockStruct{Name: "Jacek"})

		assert.False(t, act)
	})

	t.Run("given custom value string slice shouldn't find object because it not in the slice", func(t *testing.T) {
		t.Parallel()
		var Z TestingMockString = "Z"
		var W TestingMockString = "W"
		var D TestingMockString = "D"
		asset := []TestingMockString{Z, W}
		act := Contains(asset, D)

		assert.False(t, act)
	})
}

func TestSliceRemove(t *testing.T) {
	t.Run("remove item from slice", func(t *testing.T) {
		t.Parallel()
		var items = []int{1, 2, 3, 4, 5}
		var expected = []int{1, 3, 4, 5}

		var act = Remove(items, 2)

		assert.Equal(t, expected, act)

	})

	t.Run("remove items from slice", func(t *testing.T) {
		t.Parallel()
		var items = []int{1, 2, 3, 4, 5}
		var expected = []int{1}

		var act = Remove(items, 2)
		act = Remove(act, 4)
		act = Remove(act, 5)
		act = Remove(act, 3)

		assert.Equal(t, expected, act)
	})

	t.Run("remove item that doesn't exist in slice", func(t *testing.T) {
		t.Parallel()
		var items = []int{1, 2, 3, 4, 5}
		var expected = []int{1, 2, 3, 4, 5}

		var act = Remove(items, 6)

		assert.Equal(t, expected, act)
	})
}
