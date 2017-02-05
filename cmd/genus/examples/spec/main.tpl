func main() {
	app := cli.NewApp()
	app.Name = `{{ .AppName }}`
	app.Usage = `{{ .AppDescr }}`
	app.Version = `{{ .AppVersion }}`

	app.Action = func(c *cli.Context) (err error) {
		cli.ShowAppHelp(c)
		return
	}

	app.Commands = []cli.Command{serverCmd}
	app.Run(os.Args)
}

var serverFlags = []cli.Flag{
	cli.StringSliceFlag{Name: "schemes", Usage: "the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec", Value: &cli.StringSlice{"http"}},
	cli.IntFlag{Name: "cleanup-timeout", Usage: "grace period for which to wait before shutting down the server", Value: 10},
	cli.IntFlag{Name: "max-header-size", Usage: "controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body.", Value: 1024},

	cli.StringFlag{Name: "host", Usage: "the IP to listen on", Value: "localhost"},
	cli.IntFlag{Name: "port", Usage: "the port to listen on for insecure connections", Value: 8080},
	cli.IntFlag{Name: "keep-alive", Usage: "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)", Value: 60},
	cli.IntFlag{Name: "read-timeout", Usage: "maximum duration before timing out read of the request", Value: 60},
	cli.IntFlag{Name: "write-timeout", Usage: "maximum duration before timing out write of the response", Value: 60},
	cli.IntFlag{Name: "listen-limit", Usage: "limit the number of outstanding requests", Value: 10},

	cli.StringFlag{Name: "tls-host", Usage: "the IP to listen on for tls, when not specified it's the same as --host"},
	cli.IntFlag{Name: "tls-port", Usage: "the port to listen on for secure connections, when not specified it's the same as --port", Value: 8080},
	cli.StringFlag{Name: "tls-certificate", Usage: "the certificate to use for secure connections"},
	cli.StringFlag{Name: "tls-certificate-key", Usage: "the private key to use for secure conections"},
	cli.StringFlag{Name: "tls-ca-certificate", Usage: "the certificate authority file to be used with mutual tls auth"},
	cli.IntFlag{Name: "tls-keep-alive", Usage: "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)", Value: 60},
	cli.IntFlag{Name: "tls-read-timeout", Usage: "maximum duration before timing out read of the request", Value: 60},
	cli.IntFlag{Name: "tls-write-timeout", Usage: "maximum duration before timing out write of the response", Value: 60},
	cli.IntFlag{Name: "tls-listen-limit", Usage: "limit the number of outstanding requests"},

	cli.StringFlag{Name: "socket", Usage: "the unix socket to listen on", Value: "/var/run/{{ .Name }}.sock"},
}

var serverCmd = cli.Command{
	Name:      "server",
	ShortName: "s",
	Usage:     "Application Server",
	Action: func(c *cli.Context) (err error) {
		s := server.NewServer()
		if err := s.SetSchemes(c.StringSlice("schemes")...); err != nil {
			return err
		}

		if err := s.SetCleanupTimeout(c.Int("cleanup-timeout")); err != nil {
			return err
		}

		if err := s.SetMaxHeaderSize(c.Int("max-header-size")); err != nil {
			return err
		}

    s.SetHandler(server.DefaultHandler)

		if server.StringSlice(s.Schemes()).IsInclude("http") {
			if err := s.SetListenLimit(c.Int("listen-limit")); err != nil {
				return err
			}

			if err := s.SetHost(c.String("host")); err != nil {
				return err
			}

			if err := s.SetPort(c.Int("port")); err != nil {
				return err
			}

			if err := s.SetKeepAlive(c.Int("keep-alive")); err != nil {
				return err
			}

			if err := s.SetReadTimeout(c.Int("read-timeout")); err != nil {
				return err
			}

			if err := s.SetWriteTimeout(c.Int("write-timeout")); err != nil {
				return err
			}
		}

		if server.StringSlice(s.Schemes()).IsInclude("https") {
      if err := s.SetTLSHost(c.String("tls-host")); err != nil {
        return err
      }

			if err := s.SetTLSPort(c.Int("tls-port")); err != nil {
				return err
			}

			if err := s.SetTLSCertificate(c.String("tls-certificate")); err != nil {
				return err
			}

			if err := s.SetTLSCertificateKey(c.String("tls-certificate-key")); err != nil {
				return err
			}

			if err := s.SetTLSCACertificate(c.String("tls-ca-certificate")); err != nil {
				return err
			}

			if err := s.SetTLSKeepAlive(c.Int("tls-keep-alive")); err != nil {
				return err
			}

			if err := s.SetTLSReadTimeout(c.Int("tls-read-timeout")); err != nil {
				return err
			}
			if err := s.SetTLSWriteTimeout(c.Int("tls-write-timeout")); err != nil {
				return err
			}
			if err := s.SetTLSListenLimit(c.Int("tls-listen-limit")); err != nil {
				return err
			}
		}

		if server.StringSlice(s.Schemes()).IsInclude("unix") {
			if err := s.SetSocketPath(c.String("socket")); err != nil {
				return err
			}
		}

		if err := s.Serve(); err != nil {
			return err
		}

		return
	},
	Flags: serverFlags,
}
