package api_tests

import (
	"bytes"
	"github.com/donutloop/home24/internal/api"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
)

var html5 = `
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"
   "http://www.w3.org/TR/html4/strict.dtd">
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>The HTML5 Herald</title>
  <meta name="description" content="The HTML5 Herald">
  <meta name="author" content="SitePoint">
  <link rel="stylesheet" href="css/styles.css?v=1.0">
</head>
<body>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h1>Test</h1>
  <h2>Test</h2>
  <h3>Test</h3>
  <h4>Test</h4>
  <h5>Test</h5>
  <h6>Test</h6>
  <a href="/test"></a>
  <a href="http://www.example.de/test"></a>
  <script src="js/scripts.js"></script>
</body>
</html>
`

func TestValidFlow(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(html5))
	}))

	a := api.NewAPI(true)
	a.Bootstrap()
	a.Start()
	defer a.Stop()

	jsonData := `{"website_url":"` + testServer.URL + `"}`

	resp, err := http.Post(a.Server.TestURL+"/websitestats", "application/json", bytes.NewReader([]byte(jsonData)))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Fatal(err)
		}

		t.Fatal(string(respBody))
	}

	respBody, err := httputil.DumpResponse(resp, true)
	if err != nil {
		t.Fatal(err)
	}

	t.Fatal(string(respBody))
}
