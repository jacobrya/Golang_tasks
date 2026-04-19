package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both_positive", 2, 3, 5},
		{"positive_plus_zero", 5, 0, 5},
		{"negative_plus_pos", -1, 4, 3},
		{"both_negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both_positive", 5, 2, 3},
		{"positive_minus_zero", 5, 0, 5},
		{"negative_minus_positive", -2, 3, -5},
		{"both_negative", -5, -2, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("successful division", func(t *testing.T) {
		got, err := Divide(20, 4)
		require.NoError(t, err)
		require.Equal(t, 5, got)
	})

	t.Run("handling division by zero", func(t *testing.T) {
		_, err := Divide(10, 0)
		require.Error(t, err)
		require.EqualError(t, err, "division by zero")
	})
}
