var flagDebugMode = flag.Bool("test.sqldebug", false, "Turns on debug mode for SQL statements")

var (
	dbMain tester
)

type tester interface {
	setup() error
	conn() (*sql.DB, error)
	teardown() error
}

func TestMain(m *testing.M) {
	flag.Parse()
	if dbMain == nil {
		fmt.Println("no dbMain tester interface was ready")
		os.Exit(-1)
	}

	rand.Seed(time.Now().UnixNano())
	var err error

	// Load configuration
	err = initViper()
	if err != nil {
		fmt.Println("unable to load config file")
		os.Exit(-2)
	}

	setConfigDefaults()

	// Set DebugMode so we can see generated sql statements
	boil.DebugMode = *flagDebugMode

	if err = dbMain.setup(); err != nil {
		fmt.Println("Unable to execute setup:", err)
		os.Exit(-4)
	}

  conn, err := dbMain.conn()
  if err != nil {
    fmt.Println("failed to get connection:", err)
  }

	var code int
	boil.SetDB(conn)
	code = m.Run()

	if err = dbMain.teardown(); err != nil {
		fmt.Println("Unable to execute teardown:", err)
		os.Exit(-5)
	}

	os.Exit(code)
}

func initViper() error {
  var err error

	viper.SetConfigName("{{ .Schema }}")

	configHome := os.Getenv("XDG_CONFIG_HOME")
	homePath := os.Getenv("HOME")
	wd, err := os.Getwd()
	if err != nil {
		wd = "../"
	} else {
		wd = wd + "/.."
	}

	configPaths := []string{wd}
	if len(configHome) > 0 {
		configPaths = append(configPaths, filepath.Join(configHome, "{{ .Schema }}"))
	} else {
		configPaths = append(configPaths, filepath.Join(homePath, ".config/{{ .Schema }}"))
	}

	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	// Ignore errors here, fall back to defaults and validation to provide errs
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	return nil
}

// setConfigDefaults is only necessary because of bugs in viper, noted in main
func setConfigDefaults() {
	if viper.GetString("postgres.sslmode") == "" {
		viper.Set("postgres.sslmode", "require")
	}
	if viper.GetInt("postgres.port") == 0 {
		viper.Set("postgres.port", 5432)
	}
	if viper.GetString("mysql.sslmode") == "" {
		viper.Set("mysql.sslmode", "true")
	}
	if viper.GetInt("mysql.port") == 0 {
		viper.Set("mysql.port", 3306)
	}
}

