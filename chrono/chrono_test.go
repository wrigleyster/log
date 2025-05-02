package chrono

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeekday(t *testing.T) {
	monday := time.Date(2025, time.January, 6, 1, 1, 1, 1, time.UTC)
	tuesday := monday.AddDate(0, 0, 1)
	wednesday := monday.AddDate(0, 0, 2)
	thursday := monday.AddDate(0, 0, 3)
	friday := monday.AddDate(0, 0, 4)
	saturday := monday.AddDate(0, 0, 5)
	sunday := monday.AddDate(0, 0, 6)

	assert.True(t, IsWeekday(monday))
	assert.True(t, IsWeekday(tuesday))
	assert.True(t, IsWeekday(wednesday))
	assert.True(t, IsWeekday(thursday))
	assert.True(t, IsWeekday(friday))

	assert.False(t, IsWeekday(saturday))
	assert.False(t, IsWeekday(sunday))
}

