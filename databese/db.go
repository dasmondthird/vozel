package database

import (
        "fmt"
        "log"
        "os"
        "time"

        "gorm.io/driver/postgres"
        "gorm.io/gorm"
)

type User struct {
        ID         int64 `gorm:"primaryKey"`
        Username   string
        Registered bool
        VPNKey     string
        Location   string
        Expiry     time.Time
        Balance    int
}

var DB *gorm.DB

func InitDB() {
        host := os.Getenv("DB_HOST")
        port := os.Getenv("DB_PORT")
        user := os.Getenv("DB_USER")
        password := os.Getenv("DB_PASSWORD")
        dbname := os.Getenv("DB_NAME")

        dsn := fmt.Sprintf(
                "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
                host, port, user, password, dbname,
        )

        var err error
        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
                log.Fatal("Error connecting to the database:", err)
        }

        err = DB.AutoMigrate(&User{})
        if err != nil {
                log.Fatal("Error migrating the database:", err)
			}
		}
		
		func GetUserByID(userID int64) (*User, error) {
				var user User
				if err := DB.First(&user, userID).Error; err != nil {
						return nil, err
				}
				return &user, nil
		}
		
		func CreateUser(user *User) error {
				return DB.Create(user).Error
		}
		
		func UpdateUser(user *User) error {
				return DB.Save(user).Error
		}