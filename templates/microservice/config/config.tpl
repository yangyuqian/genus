// Reload config from file
func Reload(path string) (c *Config, err error) {
	c = &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = c.LoadData(data); err != nil {
		return nil, err
	}

	return c, nil
}

type Config struct {
}

// Load configuration from raw data in {{ titleCase .ConfigFormat }}
func (c *Config) LoadData(data []byte) (err error) {
	return {{ downcase .ConfigFormat }}.Unmarshal(data, c)
}
