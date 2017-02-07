var DB *gorm.DB

func SetDB(db *gorm.DB) (err error) {
  if db == nil {
    return errors.Errorf("Can not set global db to nil")
  }

  DB = db
  return
}

func GetDB() *gorm.DB {
  return DB
}
