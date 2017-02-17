{{- $defaultServer := "defaultServer" -}}

var {{ $defaultServer }} = server.NewServer()

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
	cli.DurationFlag{Name: "cleanup-timeout", Usage: "grace period for which to wait before shutting down the server", Value: time.Duration(30) * time.Second, Destination: &{{ $defaultServer }}.CleanupTimeout},
	cli.IntFlag{Name: "max-header-size", Usage: "controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line. It does not limit the size of the request body.", Value: 1024, Destination: &{{ $defaultServer }}.MaxHeaderSize},

{{ if .WithHTTP }}
	cli.StringFlag{Name: "host", Usage: "the IP to listen on", Value: "localhost", Destination: &{{ $defaultServer }}.Host},
	cli.IntFlag{Name: "port", Usage: "the port to listen on for insecure connections", Value: 8080, Destination: &{{ $defaultServer }}.Port},
	cli.DurationFlag{Name: "keep-alive", Usage: "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.KeepAlive},
	cli.DurationFlag{Name: "read-timeout", Usage: "maximum duration before timing out read of the request", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.ReadTimeout},
	cli.DurationFlag{Name: "write-timeout", Usage: "maximum duration before timing out write of the response", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.WriteTimeout},
	cli.IntFlag{Name: "listen-limit", Usage: "limit the number of outstanding requests", Value: 10, Destination: &{{ $defaultServer }}.ListenLimit},
{{ end }}

{{ if .WithHTTPS }}
	cli.StringFlag{Name: "tls-host", Usage: "the IP to listen on for tls, when not specified it's the same as --host", Destination: &{{ $defaultServer }}.TLSHost},
	cli.IntFlag{Name: "tls-port", Usage: "the port to listen on for secure connections, when not specified it's the same as --port", Value: 8080, Destination: &{{ $defaultServer }}.TLSPort},
	cli.StringFlag{Name: "tls-certificate", Usage: "the certificate to use for secure connections", Destination: &{{ $defaultServer }}.TLSCertificate},
	cli.StringFlag{Name: "tls-certificate-key", Usage: "the private key to use for secure conections", Destination: &{{ $defaultServer }}.TLSCertificateKey},
	cli.StringFlag{Name: "tls-ca-certificate", Usage: "the certificate authority file to be used with mutual tls auth", Destination: &{{ $defaultServer }}.TLSCACertificate},
	cli.DurationFlag{Name: "tls-keep-alive", Usage: "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop mid-download)", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.TLSKeepAlive},
	cli.DurationFlag{Name: "tls-read-timeout", Usage: "maximum duration before timing out read of the request", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.TLSReadTimeout},
	cli.DurationFlag{Name: "tls-write-timeout", Usage: "maximum duration before timing out write of the response", Value: time.Duration(60) * time.Second, Destination: &{{ $defaultServer }}.TLSWriteTimeout},
	cli.IntFlag{Name: "tls-listen-limit", Usage: "limit the number of outstanding requests"},
{{ end }}

{{ if .WithUnix }}
	cli.StringFlag{Name: "socket", Usage: "the unix socket to listen on", Value: "/var/run/main.sock", Destination: &{{ $defaultServer }}.SocketPath},
{{ end }}

	cli.StringFlag{Name: "config", Usage: "specify location of configuration", Value: "./config/server.yml"},
}

var serverCmd = cli.Command{
	Name:      "server",
	ShortName: "s",
	Usage:     "Application Server",
	Action: func(ctx *cli.Context) (err error) {
		if err := {{ $defaultServer }}.SetSchemes(ctx.StringSlice("schemes")...); err != nil {
			return err
		}

		{{ $defaultServer }}.SetHandler(server.DefaultHandler)

		if err := {{ $defaultServer }}.Serve(); err != nil {
			return err
		}

		return
	},
	Flags: serverFlags,
}
