package main

import (
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = int(uuid.New().ID())
	return
}

type LoanInfo struct {
    UserName  string
    BookTitle string
    LoanTime  time.Time
}

type Book struct {
	BID int `gorm:"primary_key"`
	Title string `gorm:"type:varchar(50)"`
	Author string `gorm:"type:varchar(50)"`
	PublisherYear int `gorm:"NOTNULL"`
	Genres string `gorm:"type:varchar(50)"`
	Price int `gorm:"NOTNULLBID"`
}

type Library struct{
	Name string `gorm:"type:varchar(50)"`
	LId int `gorm:"primarykey"`
	Password string `gorm:"notnull,unique"`
	IsPro bool `gorm:"default = true"`
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

type Loan struct {
	UID int 
	BID int
	Time time.Time
	User User `gorm:"foreignKey:UID;references:ID"`
	Book Book `gorm:"foreignKey:BID;references:BID"`
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
	errLO := db.AutoMigrate(&Loan{})
	if errLO != nil {
		log.Fatal("Loan not connect to database.")
	}
	fmt.Println("Loan connect to database.")
}

func addUser(db *gorm.DB){
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

func addBook(db *gorm.DB){
	var id, pub, pric int 
	var title, auth, gen string
	fmt.Scan(&id, &pub, &pric)
	fmt.Scan(&title, &auth, &gen)
	book := Book{
		BID: id,
		Title: title,
		Author: auth,
		PublisherYear: pub,
		Genres: gen,
		Price: pric,
	}
	db.Create(&book)
}

func printpos(){
	fmt.Println(`
1 : Add Book 
2 : Add User
3 : Searrch User
4 : Search Book
5 : Add Manager (need access)
6 : Add Library (need access) // not write this part yet...
7 : loan book
8 : All User
9 : All Book // not write this part yet...

0 : Close`)
}

func serchUser(db *gorm.DB){
	fmt.Println("search on user : ")
			fmt.Println(`
1 : Age
2 : Name
3 : id
4 : nid`)
	var che int
	fmt.Scanln(&che) 
	if che == 1 || che == 3{
		if che == 1 {
		fmt.Println("enter age : ")
		fmt.Scan(&che)
		user := User{}
		db.Where("Age = ?", che).Find(&user)
		fmt.Println("User is : ", user)
		}else{
		fmt.Println("enter id : ")
		fmt.Scan(&che)
		user := User{}
		db.Where("ID = ?", che).Find(&user)
		fmt.Println(user)
		}
	}
	if che == 2 || che == 4{
		if che == 2 {
		fmt.Println("enter name : ")
		var cheS string
		fmt.Scan(&cheS)
		user := User{}
		db.Where("Name = ?", cheS).Find(&user)
		fmt.Println("User is : ", user)
		}else{
		fmt.Println("enter Nid : ")
		var cheS string
		fmt.Scan(&cheS)
		user := User{}
		db.Where("Nid = ?", cheS).Find(&user)
		fmt.Println(user)
		}
	}
}

func searchBook(db *gorm.DB){
	fmt.Println("search on Book : ")
			fmt.Println(`
1 : Bid
2 : Title
3 : Publish year
4 : Price
5 : Author
6 : Genres`)
	book := []Book{}	
	var che int
	fmt.Scanln(&che) 
	if che == 1 {
		var set int
		fmt.Println("Enter Book ID : ")
		fmt.Scan(&set)
		db.Where("b_id = ?", set).First(&book)
		fmt.Println(book)
	}
	if che == 2 {
		var set string
		fmt.Println("Enter Book Title : ")
		fmt.Scan(&set)
		db.Where("title = ?", set).First(&book)
		fmt.Println(book)
	}
	if che == 3 {
		var set int
		fmt.Println("Enter Book Publish Year : ")
		fmt.Scan(&set)
		db.Where("publisher_year = ?", set).Find(&book)
		fmt.Println(book)
	}
	if che == 4 {
		var set int
		fmt.Println("Enter Book Price : ")
		fmt.Scan(&set)
		db.Where("Price = ?", set).Find(&book)
		fmt.Println(book)
	}
	if che == 5 {
		var set string
		fmt.Println("Enter Book Author : ")
		fmt.Scan(&set)
		db.Where("Author = ?", set).Find(&book)
		fmt.Println(book)
	}
	if che == 6 {
		var set string
		fmt.Println("Enter Book Generes : ")
		fmt.Scan(&set)
		db.Where("Generes = ?", set).Find(&book)
		fmt.Println(book)
	}
	
}

func checkLib(password string, db *gorm.DB) bool {
	lib := Library{}
	db.Where("password = ?", password).First(&lib)
	return password == lib.Password
}

func loanBook(bookid int, userid int, db *gorm.DB){
	user := User{}
	book := Book{}
	db.Where("id = ?", userid).First(&user)
	if user.ID == userid {
		db.Where("b_id = ?", bookid).First(&book)
		if book.BID == bookid {
			fmt.Println("wait a moment...")
			loan := Loan{
				UID: userid,
				BID: bookid,
				Time: time.Now(),
			}
			db.Create(&loan)
			fmt.Println("loan seccess")
		}else{
			fmt.Println("not have this book.")
		}
	}else {
		fmt.Println("not have this user.")
	}
}

func AllUser(db *gorm.DB) {
	user := []User{}
	var can, i int64
	db.Find(&user).Count(&can)
	for i = 0; i < can; i++{
		fmt.Print("user number : ", i + 1, " is : \n", "user id is : ", user[i].ID, "\n", "user name is : ", user[i].Name, "\n", "user age is : ", user[i].Age, "\n", "user NID is : ", user[i].Nid, "\n\n")
	}
}

func Lib(db *gorm.DB){
	for on := 1 ; on != 0; {
		printpos()
		var help int 
		fmt.Println("enter you pos")
		fmt.Scan(&help)
		if help == 0 {
			on = 0
		}
		if help == 1 {
			addBook(db)
		}
		if help == 2 {
			addUser(db)
		}
		if help == 3 {
			serchUser(db)
		}
		if help == 4 {
			searchBook(db)
		}
		if help == 5 {
			fmt.Println("Enter password : ")
			var password string 
			fmt.Scanln(&password)
			fmt.Println("wait a moment...")
			if checkLib(password, db) {

				var mid, age int
				var name string
				fmt.Println("you have access. enter new information : ")
				fmt.Scan(&mid, &age)
				fmt.Scan(&name)
				manager := Manager{
					MID: mid,
					Name: name,
					Age: age,			
				}
				db.Create(&manager)
			}else {
				fmt.Println("you dont have access")
			}
		}
		if help == 7 {
			var bookid, userid int 
			fmt.Scan(&bookid, &userid)
			loanBook(bookid, userid, db)

		}
		if help == 8 {
			AllUser(db)
		}

		if help == 9{
			BookLoans(db)
		}
	}
}

func BookLoans(db *gorm.DB) {
	var results []LoanInfo
	err := db.Table("loans"). 
	Select("users.name as user_name, books.title as book_title, loans.time as loan_time").
	Joins("join users on users.id = loans.uid").
	Joins("join books on books.b_id = loans.b_id").
	Scan(&results).Error

	if err != nil {
		fmt.Println("Error fetching loans:", err)
        return
	}

	for _, loan := range results {
        fmt.Printf("User: %s | Book: %s | Loan Time: %s\n", loan.UserName, loan.BookTitle, loan.LoanTime.Format(time.RFC3339))
    }
}
func main() {
	db := connectDB()
	Migrate(db)
	Lib(db)
}