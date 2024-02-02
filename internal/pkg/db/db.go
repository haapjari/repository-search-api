package database

// func SetupDatabase() *gorm.DB {
// 	viper.SetConfigFile(".env")
// 	viper.ReadInConfig()
//
// 	databaseUser := utils.GetDatabaseUser()
// 	databasePassword := utils.GetDatabasePassword()
// 	databaseName := utils.GetDatabaseName()
// 	databaseHost := utils.GetDatabaseHost()
// 	databasePort := utils.GetDatabasePort()
//
// 	dsn := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", databaseHost, databasePort, databaseUser, databaseName, databasePassword)
//
// 	// Open Database with ORM
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		panic("Connection to Database Failed!")
// 	}
//
// 	// Auto Creates Tables
// 	db.AutoMigrate(&models.Repository{})
// 	db.AutoMigrate(&models.Commit{})
//
// 	return db
// }
