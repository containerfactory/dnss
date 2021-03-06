package httpresolver

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"blitiri.com.ar/go/dnss/internal/dnsserver"
	"blitiri.com.ar/go/dnss/internal/testutil"

	"github.com/miekg/dns"
)

//
// === Tests ===
//

func TestSimple(t *testing.T) {
	_, ans, err := testutil.DNSQuery(DNSAddr, "test.blah.", dns.TypeA)
	if err != nil {
		t.Errorf("dns query returned error: %v", err)
	}
	if ans.(*dns.A).A.String() != "1.2.3.4" {
		t.Errorf("unexpected result: %q", ans)
	}

	_, ans, err = testutil.DNSQuery(DNSAddr, "test.blah.", dns.TypeMX)
	if err != nil {
		t.Errorf("dns query returned error: %v", err)
	}
	if ans.(*dns.MX).Mx != "mail.test.blah." {
		t.Errorf("unexpected result: %q", ans.(*dns.MX).Mx)
	}

	in, _, err := testutil.DNSQuery(DNSAddr, "unknown.", dns.TypeA)
	if err != nil {
		t.Errorf("dns query returned error: %v", err)
	}
	if in.Rcode != dns.RcodeNameError {
		t.Errorf("unexpected result: %q", in)
	}
}

//
// === Benchmarks ===
//

func BenchmarkHTTPSimple(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, _, err = testutil.DNSQuery(DNSAddr, "test.blah.", dns.TypeA)
		if err != nil {
			b.Errorf("dns query returned error: %v", err)
		}
	}
}

//
// === Test environment ===
//

// DNSHandler handles DNS-over-HTTP requests, and returns json data.
// This is used as the test server for our resolver.
func DNSHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/json")

	resp := jsonNXDOMAIN

	if r.Form["name"][0] == "test.blah." {
		switch r.Form["type"][0] {
		case "1", "A":
			resp = jsonA
		case "15", "MX":
			resp = jsonMX
		default:
			resp = jsonNXDOMAIN
		}
	}

	w.Write([]byte(resp))
}

// A record.
const jsonA = ` {
  "Status": 0, "TC": false, "RD": true, "RA": true, "AD": false, "CD": false,
  "Question": [ { "name": "test.blah.", "type": 1 }
  ],
  "Answer": [ { "name": "test.blah.", "type": 1, "TTL": 21599,
	  "data": "1.2.3.4" } ] }
`

// MX record.
const jsonMX = ` {
  "Status": 0, "TC": false, "RD": true, "RA": true, "AD": false, "CD": false,
  "Question": [ { "name": "test.blah.", "type": 15 } ],
  "Answer": [ { "name": "test.blah.", "type": 15, "TTL": 21599,
	  "data": "10 mail.test.blah." } ] }
`

// NXDOMAIN error.
const jsonNXDOMAIN = ` {
  "Status": 3, "TC": false, "RD": true, "RA": true, "AD": true, "CD": false,
  "Question": [ { "name": "doesnotexist.", "type": 15 } ],
  "Authority": [ { "name": ".", "type": 6, "TTL": 1798,
	  "data": "root. nstld. 2016052201 1800 900 604800 86400" } ] }
`

// Address where we will set up the DNS server.
var DNSAddr string

func TestMain(m *testing.M) {
	flag.Parse()

	DNSAddr = testutil.GetFreePort()

	// Test http server.
	httpsrv := httptest.NewServer(http.HandlerFunc(DNSHandler))

	// DNS to HTTPS server.
	srvURL, err := url.Parse(httpsrv.URL)
	if err != nil {
		fmt.Printf("Failed to parse test http server URL: %v\n", err)
		os.Exit(1)
	}
	r := NewJSON(srvURL, "")
	dth := dnsserver.New(DNSAddr, r, "")
	go dth.ListenAndServe()

	// Wait for the servers to start up.
	err = testutil.WaitForDNSServer(DNSAddr)
	if err != nil {
		fmt.Printf("Error waiting for the test servers to start: %v\n", err)
		fmt.Printf("Check the INFO logs for more details\n")
		os.Exit(1)
	}
	os.Exit(m.Run())
}
