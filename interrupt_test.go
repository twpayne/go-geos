package geos_test

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	"github.com/twpayne/go-geos"
)

// denseSelfCrossingLine returns a line that is very slow to buffer.
func denseSelfCrossingLine(n int) [][]float64 {
	coords := make([][]float64, 0, n)
	for i := range n {
		a := float64(i) * 0.7
		coords = append(coords, []float64{
			40 * math.Cos(a) * math.Cos(float64(i)/97),
			60 * math.Sin(a) * math.Sin(float64(i)/89),
		})
	}
	return coords
}

func skipIfNoInterruptCallback(t *testing.T) {
	t.Helper()
	if geos.VersionCompare(3, 14, 0) < 0 {
		t.Skip("GEOS < 3.14 does not support interrupt callbacks")
	}
}

func TestSetInterruptCallbackUnsupported(t *testing.T) {
	if geos.VersionCompare(3, 14, 0) >= 0 {
		t.Skip("GEOS >= 3.14 supports interrupt callbacks")
	}
	_, err := geos.NewContext().SetInterruptCallback(func() bool { return true })
	assert.IsError(t, err, errors.ErrUnsupported)
}

func TestSetInterruptCallbackStopsBuffer(t *testing.T) {
	skipIfNoInterruptCallback(t)
	c := geos.NewContext()
	line := c.NewLineString(denseSelfCrossingLine(20000))

	deadline := time.Now().Add(100 * time.Millisecond)
	clearInterrupt, err := c.SetInterruptCallback(func() bool { return time.Now().After(deadline) })
	assert.NoError(t, err)
	defer clearInterrupt()

	// go-geos panics when a GEOS operation fails, including when interrupted.
	start := time.Now()
	interrupted := func() (interrupted bool) {
		defer func() { interrupted = recover() != nil }()
		line.Buffer(1.25, 8)
		return false
	}()
	elapsed := time.Since(start)

	assert.True(t, interrupted)
	assert.True(t, elapsed < 5*time.Second)
}

func TestSetInterruptCallbackClear(t *testing.T) {
	skipIfNoInterruptCallback(t)
	c := geos.NewContext()
	line := c.NewLineString([][]float64{{0, 0}, {10, 0}, {10, 10}})
	clearInterrupt, err := c.SetInterruptCallback(func() bool { return true })
	assert.NoError(t, err)
	clearInterrupt()
	assert.NotZero(t, line.Buffer(1, 8))
}
