package main

import (
	// "fmt"
	// "time"
	// "fmt"
	"fmt"
	"log"
	// "net/http"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/gorilla/mux"
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = int(uuid.New().ID())
	return
}

type Book struct {
	Uid int `gorm:"foreignKey:User(ID)`
	BID int `gorm:"primary_key"`
	Title string `gorm:"type:varchar(50)"`
	Author string `gorm:"type:varchar(50)"`
	PublisherYear int `gorm:"NOTNULL"`
	Genres string `gorm:"type:varchar(50)"`
	Price int `gorm:"NOTNULL"`
}

type Library struct{
	Name string `gorm:"type:varchar(50)"`
	LId int `gorm:"primarykey"`
}

type User struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50)"`
	Age  int `gorm:"check:Age,Age >= 8"`
	Nid string `gorm:"varchar(10)"`
}

type Manager struct{
	MID int `gorm:"primarykey"`
	Name string `gorm:"type:varchar(50)"`
	Age  int `gorm:"check:Age, Age >= 25"`
}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=40120743 dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
  err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("User not connect to database.")
	}
	fmt.Println("User connect to database.")
	errM := db.AutoMigrate(&Manager{})
	if errM != nil {
		log.Fatal("Manager not connect to database.")
	}
	fmt.Println("Manager connect to database.")
	errL := db.AutoMigrate(&Library{})
	if errL != nil {
		log.Fatal("Library not connect to database.")
	}
	fmt.Println("Library connect to database.")
	errB := db.AutoMigrate(&Book{})
	if errB != nil {
		log.Fatal("Book not connect to database.")
	}
	fmt.Println("Book connect to database.")
}

func Lib(db *gorm.DB){
	for on := 1 ; on != 0; {
		var help int 
		fmt.Println("enter you pos")
		fmt.Scan(&help)
		if help == 0 {
			on = 0
		}
		if help == 1 {
			var id, pub, pric, uid int 
			var title, auth, gen string
			fmt.Scan(&id, &pub, &pric, &uid)
			fmt.Scan(&title, &auth, &gen)
			book := Book{
				BID: id,
				Uid: uid,
				Title: title,
				Author: auth,
				PublisherYear: pub,
				Genres: gen,
				Price: pric,
			}
			db.Create(&book)
		}
		if help == 2 {
			var id, age int 
			var name, nid string
			fmt.Scan(&id, &age)
			fmt.Scan(&nid, &name)
			user := User{
				ID: id,
				Name: name,
				Age: age,
				Nid: nid,
			}
			db.Create(&user)
		}
	}
}


func main() {
	db := connectDB()
	Migrate(db)
	Lib(db)
}