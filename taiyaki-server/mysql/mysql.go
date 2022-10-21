package mysql

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

var Db *gorm.DB

func InitDb() *gorm.DB {
	dirname, erro := os.UserHomeDir()
	if erro != nil {
		panic(erro)
	}
	err := godotenv.Load(dirname + "/.env")
	if err != nil {
		panic(err)
	}
	dbConfig := DBConfig{os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DATABASENAME")}
	Db = connectDB(dbConfig)
	return Db
}

func connectDB(c DBConfig) *gorm.DB {
	fmt.Println("Connecting to the MySQL Server")
	dsn := c.Username + ":" + c.Password + "@tcp" + "(" + c.Host + ":" + c.Port + ")/" + c.Database + "?" + "parseTime=true&loc=Local"
	fmt.Println("dsn : ", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error connecting to database : error=%v", err)
		return nil
	}
	fmt.Println("Connected to the MySQL Server")
	return db
}
