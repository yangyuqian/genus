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
	cleanupTimeout    time.Duration
	maxHeaderSize     int
	socketPath        string
	domainSocketL     net.Listener
	host              string
	port              string
	listenLimit       int
	keepAlive         time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	httpServerL       net.Listener
	tlsHost           string
	tlsPort           string
	tlsCertificate    string
	tlsCertificateKey string
	tlsCACertificate  string
	tlsListenLimit    int
	tlsKeepAlive      time.Duration
	tlsReadTimeout    time.Duration
	tlsWriteTimeout   time.Duration
	httpsServerL      net.Listener
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
		case schemeHTTP:
			s.serveHTTP(&wg)
		case schemeHTTPS:
			s.serveHTTPS(&wg)
		case schemeUnix:
			s.serveUnix(&wg)
		}
	}

	wg.Wait()
	return
}

// Serve Schemes
func (s *Server) serveHTTP(wg *sync.WaitGroup) (err error) {
	httpServer := &graceful.Server{Server: new(http.Server)}
	httpServer.MaxHeaderBytes = s.maxHeaderSize
	httpServer.SetKeepAlivesEnabled(s.keepAlive > 0)
	httpServer.TCPKeepAlive = s.keepAlive
	if s.listenLimit > 0 {
		httpServer.ListenLimit = s.listenLimit
	}

	if s.cleanupTimeout > 0 {
		httpServer.Timeout = s.cleanupTimeout
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

func (s *Server) serveHTTPS(wg *sync.WaitGroup) (err error) {
	httpsServer := &graceful.Server{Server: new(http.Server)}
	httpsServer.MaxHeaderBytes = int(s.maxHeaderSize)
	httpsServer.ReadTimeout = s.tlsReadTimeout
	httpsServer.WriteTimeout = s.tlsWriteTimeout
	httpsServer.SetKeepAlivesEnabled(s.tlsKeepAlive > 0)
	httpsServer.TCPKeepAlive = s.tlsKeepAlive
	if s.tlsListenLimit > 0 {
		httpsServer.ListenLimit = s.tlsListenLimit
	}
	if s.cleanupTimeout > 0 {
		httpsServer.Timeout = s.cleanupTimeout
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

	if s.tlsCertificate != "" && s.tlsCertificateKey != "" {
		httpsServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
		httpsServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(string(s.tlsCertificate), string(s.tlsCertificateKey))
	}

	if s.tlsCACertificate != "" {
		caCert, err := ioutil.ReadFile(string(s.tlsCACertificate))
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
		if s.tlsCertificate == "" {
			if s.tlsCertificateKey == "" {
			}
		}
		if s.tlsCertificateKey == "" {
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

func (s *Server) serveUnix(wg *sync.WaitGroup) (err error) {
	domainSocket := &graceful.Server{Server: new(http.Server)}
	domainSocket.MaxHeaderBytes = int(s.maxHeaderSize)
	domainSocket.Handler = s.handler
	if s.cleanupTimeout > 0 {
		domainSocket.Timeout = s.cleanupTimeout
	}

	wg.Add(1)
	go func(l net.Listener) {
		defer wg.Done()
		if err := domainSocket.Serve(l); err != nil {
		}
	}(s.domainSocketL)
	return
}

// Listen Scheme
func (s *Server) listen() error {
	if s.hasListeners { // already done this
		return nil
	}

	for _, scheme := range s.schemes {
		switch scheme {
		case schemeHTTP:
			if httpErr := s.listenHTTP(); httpErr != nil {
				return httpErr
			}
		case schemeHTTPS:
			if httpsErr := s.listenHTTPS(); httpsErr != nil {
				return httpsErr
			}
		case schemeUnix:
			if unixErr := s.listenUnix(); unixErr != nil {
				return unixErr
			}
		}
	}

	s.hasListeners = true
	return nil
}

func (s *Server) listenHTTP() (err error) {
	listener, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return err
	}

	s.httpServerL = listener
	return
}

func (s *Server) listenHTTPS() (err error) {
	// Use http host if https host wasn't defined
	if s.tlsHost == "" {
		s.tlsHost = s.host
	}
	// Use http listen limit if https listen limit wasn't defined
	if s.tlsListenLimit == 0 {
		s.tlsListenLimit = s.listenLimit
	}
	// Use http tcp keep alive if https tcp keep alive wasn't defined
	if s.tlsKeepAlive == 0 {
		s.tlsKeepAlive = s.keepAlive
	}
	// Use http read timeout if https read timeout wasn't defined
	if s.tlsReadTimeout == 0 {
		s.tlsReadTimeout = s.readTimeout
	}
	// Use http write timeout if https write timeout wasn't defined
	if s.tlsWriteTimeout == 0 {
		s.tlsWriteTimeout = s.writeTimeout
	}

	tlsListener, err := net.Listen("tcp", net.JoinHostPort(s.tlsHost, s.tlsPort))
	if err != nil {
		return err
	}

	s.httpsServerL = tlsListener
	return
}

func (s *Server) listenUnix() (err error) {
	domSockListener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return err
	}
	s.domainSocketL = domSockListener
	return
}

func (s *Server) Schemes() []string {
	return s.schemes
}

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

func (s *Server) SetMaxHeaderSize(maxHeaderSize int) (err error) {
	s.maxHeaderSize = maxHeaderSize
	return
}

func (s *Server) SetCleanupTimeout(cleanupTimeout int) (err error) {
	s.cleanupTimeout = time.Duration(cleanupTimeout) * time.Second
	return
}

func (s *Server) SetHost(h string) (err error) {
	s.host = h
	return
}

func (s *Server) SetPort(p int) (err error) {
	s.port = strconv.Itoa(p)
	return
}

func (s *Server) SetListenLimit(limit int) (err error) {
	s.listenLimit = limit
	return
}

func (s *Server) SetKeepAlive(d int) (err error) {
	s.keepAlive = time.Duration(d)
	return
}

func (s *Server) SetReadTimeout(timeout int) (err error) {
	s.readTimeout = time.Duration(timeout)
	return
}

func (s *Server) SetWriteTimeout(timeout int) (err error) {
	s.writeTimeout = time.Duration(timeout)
	return
}

func (s *Server) SetShutdownFunc(shutFunc shutFunc) (err error) {
	s.shutFunc = shutFunc
	return
}

func (s *Server) SetTLSCACertificate(cert string) (err error) {
	s.tlsCACertificate = cert
	return
}

func (s *Server) SetTLSCertificate(cert string) (err error) {
	s.tlsCertificate = cert
	return
}

func (s *Server) SetTLSCertificateKey(key string) (err error) {
	s.tlsCertificateKey = key
	return
}

func (s *Server) SetTLSHost(h string) (err error) {
	s.tlsHost = h
	return
}

func (s *Server) SetTLSPort(p int) (err error) {
	s.tlsPort = strconv.Itoa(p)
	return
}

func (s *Server) SetTLSKeepAlive(d int) (err error) {
	s.tlsKeepAlive = time.Duration(d)
	return
}

func (s *Server) SetTLSListenLimit(limit int) (err error) {
	s.tlsListenLimit = limit
	return
}

func (s *Server) SetTLSReadTimeout(timeout int) (err error) {
	s.tlsReadTimeout = time.Duration(timeout)
	return
}

func (s *Server) SetTLSWriteTimeout(timeout int) (err error) {
	s.tlsWriteTimeout = time.Duration(timeout)
	return
}

func (s *Server) SetSocketPath(path string) (err error) {
	s.socketPath = path
	return
}

func (s *Server) SetHandler(h http.Handler) (err error) {
	s.handler = h
	return
}
