package aggregator

import (
	// stdlib
	"testing"

	// 3p
	"github.com/stretchr/testify/assert"
)

func TestMonotonicCountSampling(t *testing.T) {
	// Initialize monotonic count
	monotonicCount := MonotonicCount{}

	// Flush w/o samples: error
	_, err := monotonicCount.flush(40)
	assert.NotNil(t, err)

	// Flush with one sample only and no prior samples: error
	monotonicCount.addSample(1, 45)
	_, err = monotonicCount.flush(40)
	assert.NotNil(t, err)

	// Add samples
	monotonicCount.addSample(2, 50)
	monotonicCount.addSample(3, 55)
	monotonicCount.addSample(6, 55)
	monotonicCount.addSample(7, 58)
	series, err := monotonicCount.flush(60)
	assert.Nil(t, err)
	if assert.Len(t, series, 1) && assert.Len(t, series[0].Points, 1) {
		assert.InEpsilon(t, 6, series[0].Points[0].Value, epsilon)
		assert.EqualValues(t, 60, series[0].Points[0].Ts)
	}

	// Flush w/o samples: error
	series, err = monotonicCount.flush(70)
	assert.NotNil(t, err)

	// Add a single sample
	monotonicCount.addSample(11, 75)
	series, err = monotonicCount.flush(80)
	assert.Nil(t, err)
	if assert.Len(t, series, 1) && assert.Len(t, series[0].Points, 1) {
		assert.InEpsilon(t, 4, series[0].Points[0].Value, epsilon)
		assert.EqualValues(t, 80, series[0].Points[0].Ts)
	}

	// Add sequence of non-monotonic samples
	monotonicCount.addSample(12, 85)
	monotonicCount.addSample(10, 85)
	monotonicCount.addSample(20, 85)
	monotonicCount.addSample(13, 85)
	monotonicCount.addSample(17, 85)
	series, err = monotonicCount.flush(90)
	assert.Nil(t, err)
	if assert.Len(t, series, 1) && assert.Len(t, series[0].Points, 1) {
		// should skip when counter is reset, i.e. between 12 and 10, and btw 20 and 13
		// 15 = (12 - 11) + (20 - 10) + (17 - 13)
		assert.InEpsilon(t, 15, series[0].Points[0].Value, epsilon)
		assert.EqualValues(t, 90, series[0].Points[0].Ts)
	}
}
