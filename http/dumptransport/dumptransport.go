package dumptransport

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

// DumpTransport is an http transport that prints the entirety of an HTTP
// roundtrip to STDOUT for debugging purposes
type DumpTransport struct {
	R http.RoundTripper
}

// RoundTrip allows DumpTransport to implement http.RoundTripper, thus it can
// act as a transport, and indeed wraps http.DefaultTransport
func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, derr := httputil.DumpRequestOut(h, true)
	if derr != nil {
		dump = []byte("!!REQUEST DUMP ERROR!! " + derr.Error())
	}

	fmt.Printf("****REQUEST****\n%s\n****RESPONSE****\n", string(dump))
	resp, err := d.R.RoundTrip(h)
	dump, derr = httputil.DumpResponse(resp, true)
	if derr != nil {
		dump = []byte("!!RESPONSE DUMP ERROR!! " + derr.Error())
	}
	fmt.Printf("%s\n****************\n\n", string(dump))
	return resp, err
}

func NewDumpTransport() *DumpTransport {
	return &DumpTransport{R: http.DefaultTransport}
}
