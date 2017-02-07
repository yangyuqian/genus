var testHTTP, testHTTPS *httptest.Server
var testServer *Server

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`OK`))
})

func TestMain(m *testing.M) {
	testServer = &Server{
  {{ if .WithHTTP }}
		ListenLimit:     10,
		CleanupTimeout:  10,
		KeepAlive:       10,
		MaxHeaderSize:   1024,
		ReadTimeout:     60,
		WriteTimeout:    60,
  {{ end }}
  {{ if .WithHTTPS }}
		TLSKeepAlive:    10,
		TLSReadTimeout:  60,
		TLSWriteTimeout: 60,
		TLSListenLimit:  10,
  {{ end }}
		hasListeners:    true,
	}
  {{ if .WithHTTP }}
	testHTTP = httptest.NewServer(nil)
	testServer.httpServerL = testHTTP.Listener
  {{ end }}

  {{ if .WithHTTPS }}
	testHTTPS = httptest.NewTLSServer(nil)
	testServer.httpsServerL = testHTTPS.Listener
  {{ end }}

  {{ if .WithUnix }}
	testServer.domainSocketL = testHTTP.Listener
  {{ end }}

	testServer.SetHandler(testHandler)

	code := m.Run()
	os.Exit(code)
}

func TestServer(t *testing.T) {
	t.Run("Serve", testServe)
  {{ if .WithHTTP }}
	t.Run("ServeHTTP", testServeHTTP)
  {{ end }}
  {{ if .WithHTTPS }}
	t.Run("ServeHTTPS", testServeHTTPS)
  {{ end }}
  {{ if .WithUnix }}
	t.Run("ServeUnix", testServeUnix)
  {{ end }}
}

func testServe(t *testing.T) {
	t.Parallel()
	if err := testServer.Serve(); err != nil {
		t.Errorf("Serve error <%v>", err)
	}
}

{{ if .WithHTTP }}
func testServeHTTP(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	testServer.serveHTTP(wg)
}
{{ end }}

{{ if .WithHTTPS }}
func testServeHTTPS(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	testServer.serveHTTPS(wg)
}
{{ end }}

{{ if .WithUnix }}
func testServeUnix(t *testing.T) {
	t.Parallel()
	wg := &sync.WaitGroup{}
	if err := testServer.serveUnix(wg); err != nil {
		t.Errorf("serve on socket error <%+v>", err)
	}
}
{{ end }}
