package dumptransport

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)

// DumpTransport is an http transport that prints the entirety of an HTTP
// roundtrip to STDOUT for debugging purposes
type DumpTransport struct {
	r http.RoundTripper
}

// RoundTrip allows DumpTransport to implement http.RoundTripper, thus it can
// act as a transport, and indeed wraps http.DefaultTransport
func (d *DumpTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	dump, derr := httputil.DumpRequestOut(h, true)
	if derr != nil {
		dump = []byte("!!REQUEST DUMP ERROR!! " + derr.Error())
	}

	fmt.Printf("****REQUEST****\n%s\n****RESPONSE****\n", dump2Str(dump))
	resp, err := d.r.RoundTrip(h)
	dump, derr = httputil.DumpResponse(resp, true)
	if derr != nil {
		dump = []byte("!!RESPONSE DUMP ERROR!! " + derr.Error())
	}
	fmt.Printf("%s\n****************\n\n", dump2Str(dump))
	return resp, err
}

func NewDumpTransport() *DumpTransport {
	return &DumpTransport{r: http.DefaultTransport}
}

func dump2Str(b []byte) string {
	// excape string, then add actual newlines to escaped newlines for display
	// purposes
	qStr := strings.Replace(
		strconv.QuoteToASCII(string(b)),
		"\\n",
		"\\n\n",
		-1)
	// strconv.QuoteToASCII adds quotes to string, remove final quote
	idx := len(qStr) - 1
	// if last character is a newline, don't add an extra one
	if qStr[idx-1] == '\n' {
		idx--
	}
	// remove opening quote
	return qStr[1:idx]
}
