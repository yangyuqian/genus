func TestMain(m *testing.M) {
	flag.Parse()

  tester := &mysqlTester{}
  tester.setup()
  defer tester.teardown()

  conn, _ := tester.conn()
  
  SetDB(conn)

	code := m.Run()
  os.Exit(code)
}
