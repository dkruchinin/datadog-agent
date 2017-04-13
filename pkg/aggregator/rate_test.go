package aggregator

import (
	// stdlib
	"testing"

	// 3p
	"github.com/stretchr/testify/assert"
)

func TestRateSampling(t *testing.T) {
	// Initialize rates
	mRate1 := Rate{}
	mRate2 := Rate{}

	// Add samples
	mRate1.addSample(1, 50)
	mRate1.addSample(2, 55)
	mRate2.addSample(1, 50)

	// First rate
	series, err := mRate1.flush(60)
	assert.Nil(t, err)
	assert.Len(t, series, 1)
	assert.Len(t, series[0].Points, 1)
	assert.InEpsilon(t, 0.2, series[0].Points[0].Value, epsilon)
	assert.EqualValues(t, 55, series[0].Points[0].Ts)

	// Second rate (should return error)
	_, err = mRate2.flush(60)
	assert.NotNil(t, err)
}

func TestRateSamplingMultipleSamplesInSameFlush(t *testing.T) {
	// Initialize rate
	mRate := Rate{}

	// Add samples
	mRate.addSample(1, 50)
	mRate.addSample(2, 55)
	mRate.addSample(4, 61)

	// Should compute rate based on the last 2 samples
	series, err := mRate.flush(65)
	assert.Nil(t, err)
	assert.Len(t, series, 1)
	assert.Len(t, series[0].Points, 1)
	assert.InEpsilon(t, 2./6., series[0].Points[0].Value, epsilon)
	assert.EqualValues(t, 61, series[0].Points[0].Ts)
}

func TestRateSamplingNoSampleForOneFlush(t *testing.T) {
	// Initialize rate
	mRate := Rate{}

	// Add samples
	mRate.addSample(1, 50)
	mRate.addSample(2, 55)

	// First flush: no error
	_, err := mRate.flush(60)
	assert.Nil(t, err)

	// Second flush w/o sample: error
	_, err = mRate.flush(60)
	assert.NotNil(t, err)

	// Third flush w/ sample
	mRate.addSample(4, 60)
	// Should compute rate based on the last 2 samples
	series, err := mRate.flush(60)
	assert.Nil(t, err)
	assert.Len(t, series, 1)
	assert.Len(t, series[0].Points, 1)
	assert.InEpsilon(t, 2./5., series[0].Points[0].Value, epsilon)
	assert.EqualValues(t, 60, series[0].Points[0].Ts)
}

func TestRateSamplingSamplesAtSameTimestamp(t *testing.T) {
	// Initialize rate
	mRate := Rate{}

	// Add samples
	mRate.addSample(1, 50)
	mRate.addSample(2, 50)

	series, err := mRate.flush(60)

	assert.NotNil(t, err)
	assert.Len(t, series, 0)
}
