package nettools

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddress(t *testing.T) {
	for i, tt := range []struct {
		addr    string
		resAddr *Address
		failed  bool
	}{
		{addr: "127.0.0.1:5000", resAddr: &Address{Host: "127.0.0.1", Port: 5000}},
		{addr: ":5000", resAddr: &Address{Host: "", Port: 5000}},
		{addr: "127.0.0.1", failed: true},
		{addr: "127.0.0.1:", failed: true},
		{addr: "5000", failed: true},
		{addr: "127.0.0.1:aaaa", failed: true},
	} {
		i, tt := i, tt
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			resAddr, err := NewAddress(tt.addr)
			if tt.failed {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.resAddr, resAddr)
				assert.Equal(t, tt.addr, resAddr.String())
			}
		})
	}
}
