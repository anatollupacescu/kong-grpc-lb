package product

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	date := time.Now()

	tt := []struct {
		name     string
		products []CSVRow
		expected []productDTO
	}{
		{
			name: "single product, one update",
			products: []CSVRow{
				{Name: "test", Price: "2.0"},
				{Name: "test", Price: "1.0"},
			},
			expected: []productDTO{{
				LastUpdated: date,
				UpdateCount: 1,
				Name:        "test",
				Price:       "1.0",
			}},
		}, {
			name: "single product, two updates",
			products: []CSVRow{
				{Name: "test", Price: "2.0"},
				{Name: "test", Price: "1.0"},
				{Name: "test", Price: "9.0"},
			},
			expected: []productDTO{{
				LastUpdated: date,
				UpdateCount: 2,
				Name:        "test",
				Price:       "9.0",
			}},
		}, {
			name: "two products",
			products: []CSVRow{
				{Name: "test 2", Price: "2.0"},
				{Name: "test", Price: "1.0"},
				{Name: "test 2", Price: "9.0"},
			},
			expected: []productDTO{
				{
					LastUpdated: date,
					UpdateCount: 1,
					Name:        "test 2",
					Price:       "9.0",
				},
				{
					LastUpdated: date,
					UpdateCount: 0,
					Name:        "test",
					Price:       "1.0",
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			res := transform(test.products, date)
			assert.Equal(t, test.expected, res)
		})
	}
}
