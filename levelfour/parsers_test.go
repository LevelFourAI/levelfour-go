package levelfour

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		got, err := ParseDate("2024-03-15")
		require.NoError(t, err)
		assert.Equal(t, 2024, got.Year())
		assert.Equal(t, time.March, got.Month())
		assert.Equal(t, 15, got.Day())
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := ParseDate("not-a-date")
		assert.Error(t, err)
	})

	t.Run("wrong format", func(t *testing.T) {
		_, err := ParseDate("15/03/2024")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := ParseDate("")
		assert.Error(t, err)
	})
}

func TestParseDateTime(t *testing.T) {
	t.Run("valid datetime", func(t *testing.T) {
		got, err := ParseDateTime("2024-03-15T10:30:00Z")
		require.NoError(t, err)
		assert.Equal(t, 2024, got.Year())
		assert.Equal(t, time.March, got.Month())
		assert.Equal(t, 10, got.Hour())
	})

	t.Run("with timezone offset", func(t *testing.T) {
		got, err := ParseDateTime("2024-03-15T10:30:00-05:00")
		require.NoError(t, err)
		assert.Equal(t, 10, got.Hour())
	})

	t.Run("invalid datetime", func(t *testing.T) {
		_, err := ParseDateTime("not-a-datetime")
		assert.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		_, err := ParseDateTime("")
		assert.Error(t, err)
	})
}

func TestParseDatePtr(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result, err := ParseDatePtr(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("valid date", func(t *testing.T) {
		s := "2024-03-15"
		result, err := ParseDatePtr(&s)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2024, result.Year())
	})

	t.Run("invalid date", func(t *testing.T) {
		s := "bad"
		_, err := ParseDatePtr(&s)
		assert.Error(t, err)
	})
}

func TestParseDateTimePtr(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result, err := ParseDateTimePtr(nil)
		require.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("valid datetime", func(t *testing.T) {
		s := "2024-03-15T10:30:00Z"
		result, err := ParseDateTimePtr(&s)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 10, result.Hour())
	})

	t.Run("invalid datetime", func(t *testing.T) {
		s := "bad"
		_, err := ParseDateTimePtr(&s)
		assert.Error(t, err)
	})
}
