const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
	schemeUnix  = "unix"
)

type shutFunc func() error

func NewServer() (s *Server) {
	return &Server{}
}

type Server struct {
	schemes           []string // HTTPS, HTTP, UNIX Socket
	CleanupTimeout    time.Duration
	MaxHeaderSize     int

{{ if .WithHTTP }}
	Host              string
	Port              int
	ListenLimit       int
	KeepAlive         time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	httpServerL       net.Listener
{{ end }}

{{ if .WithHTTPS }}
	TLSHost           string
	TLSPort           int
	TLSCertificate    string
	TLSCertificateKey string
	TLSCACertificate  string
	TLSListenLimit    int
	TLSKeepAlive      time.Duration
	TLSReadTimeout    time.Duration
	TLSWriteTimeout   time.Duration
	httpsServerL      net.Listener
{{ end }}

{{ if .WithUnix }}
	SocketPath        string
	domainSocketL     net.Listener
{{ end }}

	handler           http.Handler
	hasListeners      bool
	shutFunc          shutFunc
}

func (s *Server) Serve() (err error) {
	if s.shutFunc != nil {
		defer s.shutFunc()
	}

	if listenErr := s.listen(); listenErr != nil {
		return listenErr
	}

	var wg sync.WaitGroup

	for _, scheme := range s.schemes {
		switch scheme {
    {{ if .WithHTTP }}
		case schemeHTTP:
			s.serveHTTP(&wg)
    {{ end }}
    {{ if .WithHTTPS }}
		case schemeHTTPS:
			s.serveHTTPS(&wg)
    {{ end }}
    {{ if .WithUnix }}
		case schemeUnix:
			s.serveUnix(&wg)
    {{ end }}
		}
	}

	wg.Wait()
	return
}

{{ if .WithHTTP }}
// Serve Schemes
func (s *Server) serveHTTP(wg *sync.WaitGroup) (err error) {
	httpServer := &graceful.Server{Server: new(http.Server)}
	httpServer.MaxHeaderBytes = s.MaxHeaderSize
	httpServer.SetKeepAlivesEnabled(s.KeepAlive > 0)
	httpServer.TCPKeepAlive = s.KeepAlive * time.Second
	if s.ListenLimit > 0 {
		httpServer.ListenLimit = s.ListenLimit
	}

	if s.CleanupTimeout > 0 {
		httpServer.Timeout = s.CleanupTimeout * time.Second
	}

	httpServer.Handler = s.handler

	wg.Add(1)
	go func(l net.Listener) {
		defer wg.Done()
		if err := httpServer.Serve(l); err != nil {
		}
	}(s.httpServerL)
	return
}
{{ end }}

{{ if .WithHTTPS }}
func (s *Server) serveHTTPS(wg *sync.WaitGroup) (err error) {
	httpsServer := &graceful.Server{Server: new(http.Server)}
	httpsServer.MaxHeaderBytes = int(s.MaxHeaderSize)
	httpsServer.ReadTimeout = s.TLSReadTimeout * time.Second
	httpsServer.WriteTimeout = s.TLSWriteTimeout * time.Second
	httpsServer.SetKeepAlivesEnabled(s.TLSKeepAlive > 0)
	httpsServer.TCPKeepAlive = s.TLSKeepAlive * time.Second
	if s.TLSListenLimit > 0 {
		httpsServer.ListenLimit = s.TLSListenLimit
	}
	if s.CleanupTimeout > 0 {
		httpsServer.Timeout = s.CleanupTimeout * time.Second
	}
	httpsServer.Handler = s.handler

	// Inspired by https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
	httpsServer.TLSConfig = &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		// https://github.com/golang/go/tree/master/src/crypto/elliptic
		CurvePreferences: []tls.CurveID{tls.CurveP256},
		// Use modern tls mode https://wiki.mozilla.org/Security/Server_Side_TLS#Modern_compatibility
		NextProtos: []string{"http/1.1", "h2"},
		// https://www.owasp.org/index.php/Transport_Layer_Protection_Cheat_Sheet#Rule_-_Only_Support_Strong_Protocols
		MinVersion: tls.VersionTLS12,
		// These ciphersuites support Forward Secrecy: https://en.wikipedia.org/wiki/Forward_secrecy
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	if s.TLSCertificate != "" && s.TLSCertificateKey != "" {
		httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
		httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.TLSCertificate), string(s.TLSCertificateKey))
	}

	if s.TLSCACertificate != "" {
		caCert, err := ioutil.ReadFile(string(s.TLSCACertificate))
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		httpsServer.TLSConfig.ClientCAs = caCertPool
		httpsServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	httpsServer.TLSConfig.BuildNameToCertificate()

	if err != nil {
		return err
	}

	if len(httpsServer.TLSConfig.Certificates) == 0 {
		if s.TLSCertificate == "" {
			if s.TLSCertificateKey == "" {
			}
		}
		if s.TLSCertificateKey == "" {
		}
	}

	wg.Add(1)
	go func(l net.Listener) {
		defer wg.Done()
		if err := httpsServer.Serve(l); err != nil {
		}
	}(tls.NewListener(s.httpsServerL, httpsServer.TLSConfig))
	return
}
{{ end }}

{{ if .WithUnix }}
func (s *Server) serveUnix(wg *sync.WaitGroup) (err error) {
	domainSocket := &graceful.Server{Server: new(http.Server)}
	domainSocket.MaxHeaderBytes = int(s.MaxHeaderSize)
	domainSocket.Handler = s.handler
	if s.CleanupTimeout > 0 {
		domainSocket.Timeout = s.CleanupTimeout * time.Second
	}

	wg.Add(1)
	go func(l net.Listener) {
		defer wg.Done()
		if err := domainSocket.Serve(l); err != nil {
		}
	}(s.domainSocketL)
	return
}
{{ end }}

// Listen Scheme
func (s *Server) listen() error {
	if s.hasListeners { // already done this
		return nil
	}

	for _, scheme := range s.schemes {
		switch scheme {
    {{ if .WithHTTP }}
		case schemeHTTP:
			if httpErr := s.listenHTTP(); httpErr != nil {
				return httpErr
			}
    {{ end }}
    {{ if .WithHTTPS }}
		case schemeHTTPS:
			if httpsErr := s.listenHTTPS(); httpsErr != nil {
				return httpsErr
			}
    {{ end }}
    {{ if .WithUnix }}
		case schemeUnix:
			if unixErr := s.listenUnix(); unixErr != nil {
				return unixErr
			}
    {{ end }}
    default:
      return errors.Errorf("Unkown scheme: %+v", scheme)
		}
	}

	s.hasListeners = true
	return nil
}

{{ if .WithHTTP }}
func (s *Server) listenHTTP() (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		return err
	}

	s.httpServerL = listener
	return
}
{{ end }}

{{ if .WithHTTPS }}
func (s *Server) listenHTTPS() (err error) {
	// Use http host if https host wasn't defined
	if s.TLSHost == "" {
		s.TLSHost = s.Host
	}
	// Use http listen limit if https listen limit wasn't defined
	if s.TLSListenLimit == 0 {
		s.TLSListenLimit = s.ListenLimit
	}
	// Use http tcp keep alive if https tcp keep alive wasn't defined
	if s.TLSKeepAlive == 0 {
		s.TLSKeepAlive = s.KeepAlive
	}
	// Use http read timeout if https read timeout wasn't defined
	if s.TLSReadTimeout == 0 {
		s.TLSReadTimeout = s.ReadTimeout
	}
	// Use http write timeout if https write timeout wasn't defined
	if s.TLSWriteTimeout == 0 {
		s.TLSWriteTimeout = s.WriteTimeout
	}

	if s.TLSPort <= 0 {
		s.TLSPort = s.Port
	}

	tlsListener, err := net.Listen("tcp", net.JoinHostPort(s.TLSHost, strconv.Itoa(s.TLSPort)))
	if err != nil {
		return err
	}

	s.httpsServerL = tlsListener
	return
}
{{ end }}

{{ if .WithUnix }}
func (s *Server) listenUnix() (err error) {
	domSockListener, err := net.Listen("unix", s.SocketPath)
	if err != nil {
		return err
	}
	s.domainSocketL = domSockListener
	return
}
{{ end }}

// Setters
func (s *Server) SetSchemes(names ...string) (err error) {
	for _, name := range names {
		scheme := strings.ToLower(name)
		if scheme != "http" && scheme != "https" && scheme != "unix" {
			return errors.New(fmt.Sprintf("Unkown scheme %s", scheme))
		}

		s.schemes = append(s.schemes, scheme)
	}

	return
}

func (s *Server) SetShutdownFunc(shutFunc shutFunc) (err error) {
	s.shutFunc = shutFunc
	return
}

func (s *Server) SetHandler(h http.Handler) (err error) {
	s.handler = h
	return
}
