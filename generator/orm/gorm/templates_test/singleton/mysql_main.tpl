type mysqlTester struct {
	dbConn *gorm.DB

	dbName	string
	host	string
	user	string
	pass	string
	sslmode	string
	port	int

	optionFile string

	testDBName string
}

var dbname, host, user, pass, sslmode string
var port int

func init() {
	flag.StringVar(&dbname, "test.mysql.dbname", "", "Set test MySQL database name")
	flag.StringVar(&host, "test.mysql.host", "127.0.0.1", "Set test MySQL database host")
	flag.StringVar(&user, "test.mysql.user", "root", "Set test MySQL database username")
	flag.StringVar(&pass, "test.mysql.pass", "root", "Set test MySQL database password")
	flag.StringVar(&sslmode, "test.mysql.sslmode", "false", "Set test MySQL database sslmode")
	flag.IntVar(&port, "test.mysql.port", 3306, "Set test MySQL database port")
}

func (m *mysqlTester) setup() error {
	var err error

	m.dbName = dbname
	m.host = host
	m.user = user
	m.pass = pass
	m.port = port
	m.sslmode = sslmode
	// Create a randomized db name.
	m.testDBName = randomize.StableDBName(m.dbName)

	if err = m.dropTestDB(); err != nil {
		return err
	}
	if err = m.createTestDB(); err != nil {
		return err
	}

	dumpCmd := exec.Command("mysqldump", "-h"+m.host, "-u"+m.user, "-p"+m.pass, fmt.Sprintf("-P%d", m.port), "--no-data", m.dbName)
	createCmd := exec.Command("mysql", "-h"+m.host, "-u"+m.user, "-p"+m.pass, fmt.Sprintf("-P%d", m.port), "--database", m.testDBName)

	r, w := io.Pipe()
	dumpCmd.Stdout = w
	createCmd.Stdin = newFKeyDestroyer(rgxMySQLkey, r)

	if err = dumpCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start mysqldump command")
	}
	if err = createCmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start mysql command")
	}

	if err = dumpCmd.Wait(); err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "failed to wait for mysqldump command")
	}

	w.Close()

	if err = createCmd.Wait(); err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "failed to wait for mysql command")
	}

	return nil
}

func (m *mysqlTester) sslMode(mode string) string {
	switch mode {
	case "true":
		return "REQUIRED"
	case "false":
		return "DISABLED"
	default:
		return "PREFERRED"
	}
}

func (m *mysqlTester) defaultsFile() string {
	return fmt.Sprintf("--defaults-file=%s", m.optionFile)
}

func (m *mysqlTester) makeOptionFile() error {
	tmp, err := ioutil.TempFile("", "optionfile")
	if err != nil {
		return errors.Wrap(err, "failed to create option file")
	}

	fmt.Fprintln(tmp, "[client]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	fmt.Fprintf(tmp, "password=%s\n", m.pass)
	fmt.Fprintf(tmp, "ssl-mode=%s\n", m.sslMode(m.sslmode))

	fmt.Fprintln(tmp, "[mysqldump]")
	fmt.Fprintf(tmp, "host=%s\n", m.host)
	fmt.Fprintf(tmp, "port=%d\n", m.port)
	fmt.Fprintf(tmp, "user=%s\n", m.user)
	fmt.Fprintf(tmp, "password=%s\n", m.pass)
	fmt.Fprintf(tmp, "ssl-mode=%s\n", m.sslMode(m.sslmode))

	m.optionFile = tmp.Name()

	return tmp.Close()
}

func (m *mysqlTester) createTestDB() error {
	sql := fmt.Sprintf("create database %s;", m.testDBName)
	return m.runCmd(sql, "mysql", "-h"+m.host, "-u"+m.user, "-p"+m.pass, fmt.Sprintf("-P%d", m.port))
}

func (m *mysqlTester) dropTestDB() error {
	sql := fmt.Sprintf("drop database if exists %s;", m.testDBName)
	return m.runCmd(sql, "mysql", "-h"+m.host, "-u"+m.user, "-p"+m.pass, fmt.Sprintf("-P%d", m.port))
}

func (m *mysqlTester) teardown() error {
	if m.dbConn != nil {
		m.dbConn.Close()
	}

	if err := m.dropTestDB(); err != nil {
		return err
	}

  return nil
}

func (m *mysqlTester) runCmd(stdin, command string, args ...string) error {
	// args = append([]string{m.defaultsFile()}, args...)

	cmd := exec.Command(command, args...)
	cmd.Stdin = strings.NewReader(stdin)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
	fmt.Println("failed running:", command, args)
	fmt.Println(stdout.String())
	fmt.Println(stderr.String())
	return err
	}

	return nil
}

func (m *mysqlTester) conn() (*gorm.DB, error) {
	if m.dbConn != nil {
	return m.dbConn, nil
	}

	var err error
	m.dbConn, err = gorm.Open("mysql", drivers.MySQLBuildQueryString(m.user, m.pass, m.testDBName, m.host, m.port, m.sslmode))
	if err != nil {
	return nil, err
	}

	return m.dbConn, nil
}
