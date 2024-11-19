package clientip

import (
	"net/http"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseHTTPRequest(t *testing.T) {
	t.Parallel()

	makeHeader := func(keyValues map[string][]string) http.Header {
		header := http.Header{}
		for key, values := range keyValues {
			for _, value := range values {
				header.Add(key, value)
			}
		}
		return header
	}

	testCases := map[string]struct {
		request    *http.Request
		ip         netip.Addr
		errWrapped error
		errMessage string
	}{
		"nil request": {
			errWrapped: ErrRequestIsNil,
			errMessage: "request is nil",
		},
		"empty request": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
			},
			ip: netip.AddrFrom4([4]byte{99, 99, 99, 99}),
		},
		"request with remote address": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
			},
			ip: netip.AddrFrom4([4]byte{99, 99, 99, 99}),
		},
		"request with xRealIP header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with xRealIP header and public XForwardedFor IP": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"77.77.77.77"},
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with xRealIP header and private XForwardedFor IP": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"88.88.88.88"},
					"X-Forwarded-For": {"10.0.0.5"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with single public IP in xForwardedFor header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with two public IPs in xForwardedFor header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88", "77.77.77.77"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with private and public IPs in xForwardedFor header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "88.88.88.88", "10.0.0.1", "77.77.77.77"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{88, 88, 88, 88}),
		},
		"request with single private IP in xForwardedFor header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{192, 168, 1, 5}),
		},
		"request with private IPs in xForwardedFor header": {
			request: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "10.0.0.17"},
				}),
			},
			ip: netip.AddrFrom4([4]byte{192, 168, 1, 5}),
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ip, err := ParseHTTPRequest(testCase.request)
			assert.Equal(t, testCase.ip, ip)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
