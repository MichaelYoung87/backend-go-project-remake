package database
import (
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/MichaelYoung87/backend-go-project-remake/models"
)
var DB *gorm.DB
func Connect() {
	connection, err := gorm.Open(mysql.Open("username:password@tcp(XX.XXX.XX.XXX:3306)/NAMN_PÅ_DATABAS?charset=utf8mb4&parseTime=True&loc=Europe%2FStockholm"), &gorm.Config{})
	if err != nil {
		panic("could not connect to the server")
	}
	DB = connection
	connection.AutoMigrate(&models.User{})
}