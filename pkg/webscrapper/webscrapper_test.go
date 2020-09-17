package webscrapper

import (
	"net/http"
	"net/http/httptest"
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
  <section id="site-signin" class="is-centered-button">
  </section>
</body>
</html>
`

func TestExtract(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(html5))
	}))

	defer testServer.Close()

	webScrapper := New(new(http.Client))

	websiteData, err := webScrapper.Extract(testServer.URL)
	if err != nil {
		t.Fatal(err)
	}

	doctype := "-//W3C//DTD HTML 4.01//EN"
	if websiteData.HTMLVersion != doctype {
		t.Errorf("HTMLVersion is bad, got: %v ,want: %v", websiteData.HTMLVersion, doctype)
	}

	HeadingsCountLevel1Want := 8
	if websiteData.HeadingsCountLevel1 != HeadingsCountLevel1Want {
		t.Errorf("HeadingsCountLevel1 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel1, HeadingsCountLevel1Want)
	}

	HeadingsCountLevel2Want := 1
	if websiteData.HeadingsCountLevel2 != HeadingsCountLevel2Want {
		t.Errorf("HeadingsCountLevel3 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel2, HeadingsCountLevel2Want)
	}

	HeadingsCountLevel3Want := 1
	if websiteData.HeadingsCountLevel3 != HeadingsCountLevel3Want {
		t.Errorf("HeadingsCountLevel1 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel3, HeadingsCountLevel3Want)
	}

	HeadingsCountLevel4Want := 1
	if websiteData.HeadingsCountLevel4 != HeadingsCountLevel4Want {
		t.Errorf("HeadingsCountLevel1 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel4, HeadingsCountLevel4Want)
	}

	HeadingsCountLevel5Want := 1
	if websiteData.HeadingsCountLevel5 != HeadingsCountLevel5Want {
		t.Errorf("HeadingsCountLevel1 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel5, HeadingsCountLevel5Want)
	}

	HeadingsCountLevel6Want := 1
	if websiteData.HeadingsCountLevel6 != HeadingsCountLevel6Want {
		t.Errorf("HeadingsCountLevel1 is bad, got: %v ,want: %v", websiteData.HeadingsCountLevel6, HeadingsCountLevel6Want)
	}

	internalLinkCount := 1
	if websiteData.InternalLinkCount != internalLinkCount {
		t.Errorf("InternalLinkCount is bad, got: %v ,want: %v", websiteData.InternalLinkCount, internalLinkCount)
	}

	externalLinkCount := 1
	if websiteData.ExternalLinkCount != externalLinkCount {
		t.Errorf("ExternalLinkCount is bad, got: %v ,want: %v", websiteData.ExternalLinkCount, externalLinkCount)
	}

	PageTitle := "The HTML5 Herald"
	if websiteData.PageTitle != PageTitle {
		t.Errorf("PageTitle is bad, got: %v ,want: %v", websiteData.PageTitle, PageTitle)
	}

	hasLogin := true
	if websiteData.HasLoginForm != hasLogin {
		t.Errorf("HasLoginForm is bad, got: %v ,want: %v", websiteData.HasLoginForm, hasLogin)
	}
}
