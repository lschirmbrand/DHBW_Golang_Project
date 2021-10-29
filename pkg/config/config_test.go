package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test if all default values are set
func TestConfigure(t *testing.T) {

	Configure()

	assert.Equal(t, 8443, *CheckinPort)
	assert.Equal(t, 8444, *QRCodePort)
	assert.Equal(t, 60, *RefreshTime)
}
