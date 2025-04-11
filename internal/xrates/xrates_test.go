package xrates

import (
	"testing"

	"github.com/srackham/go-utils/assert"
	"github.com/srackham/go-utils/helpers"
)

func TestExchangeRates(t *testing.T) {
	if helpers.IsRunningOnGithub() {
		t.Skip("skip on Github Actions because this test requires HTTP access")
	}

	// The cache is not loaded or saved by this test to force it to go and fetch the exchange rates
	x := New()
	x.ConfigFile = "../../testdata/config.yaml"

	rate, err := x.GetCachedRate("USD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.Equal(t, 1.00, rate)

	rate, err = x.GetCachedRate("NZD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, rate > 0, "invalid NZD rate: %f", rate)

	rate, err = x.GetCachedRate("AUD", false)
	assert.PassIf(t, err == nil, "%v", err)
	assert.PassIf(t, rate > 0, "invalid AUD rate: %f", rate)

	_, err = x.GetCachedRate("", false)
	assert.PassIf(t, err != nil, "should have return error for blank currency")
	assert.Equal(t, "no currency specified", err.Error())

	_, err = x.GetCachedRate("FOOBAR", false)
	assert.PassIf(t, err != nil, "should have return error for FOOBAR currency")
	assert.Equal(t, "unknown currency: FOOBAR", err.Error())
}
