package generators

import (
	"math/rand"
)

const (
	variousCount = 3
	minLen       = 3
)

func GeneratePlenty[T any](count int, factory func() T) []T {
	items := make([]T, 0, count)
	for range count {
		items = append(items, factory())
	}
	return items
}

func RandomCount() int {
	return rand.Intn(variousCount) + minLen //nolint:gosec // Ok for tests.
}

func RandomInt(minV, maxV int) int {
	return rand.Intn(maxV-minV) + minV //nolint:gosec // Ok for tests.
}

func set[T any](field *T, generateValue func() T, value ...T) {
	if len(value) == 0 {
		value = []T{generateValue()}
	}

	*field = value[0]
}
