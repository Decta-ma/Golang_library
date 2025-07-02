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
	BID           int    `gorm:"primary_key"`
	Title         string `gorm:"type:varchar(50)"`
	Author        string `gorm:"type:varchar(50)"`
	PublisherYear int    `gorm:"NOTNULL"`
	Genres        string `gorm:"type:varchar(50)"`
	Price         int    `gorm:"NOTNULLBID"`
}

type Library struct {
	Name     string `gorm:"type:varchar(50)"`
	LId      int    `gorm:"primarykey"`
	Password string `gorm:"notnull,unique"`
	IsPro    bool   `gorm:"default = true"`
}

type User struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50)"`
	Age  int    `gorm:"check:Age,Age >= 8"`
	Nid  string `gorm:"varchar(10)"`
}

type Manager struct {
	MID  int    `gorm:"primarykey"`
	Name string `gorm:"type:varchar(50)"`
	Age  int    `gorm:"check:Age, Age >= 25"`
}

type Loan struct {
	UID  int
	BID  int
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

func AddManager(db *gorm.DB) {
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
			MID:  mid,
			Name: name,
			Age:  age,
		}

		db.Create(&manager)
	} else {
		fmt.Println("you dont have access")
	}
}

func addUser(db *gorm.DB) {
	var id, age int
	var name, nid string
	fmt.Scan(&id, &age)
	fmt.Scan(&nid, &name)
	user := User{
		ID:   id,
		Name: name,
		Age:  age,
		Nid:  nid,
	}
	db.Create(&user)
}

func addBook(db *gorm.DB) {
	var id, pub, pric int
	var title, auth, gen string
	fmt.Scan(&id, &pub, &pric)
	fmt.Scan(&title, &auth, &gen)
	book := Book{
		BID:           id,
		Title:         title,
		Author:        auth,
		PublisherYear: pub,
		Genres:        gen,
		Price:         pric,
	}
	db.Create(&book)
}

func printpos() {
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
10 : show user name and id in loan 
11 : update user

0 : Close`)
}

func serchUser(db *gorm.DB) {
	fmt.Println("search on user : ")
	fmt.Println(`
1 : Age
2 : Name
3 : id
4 : nid`)
	var che int
	fmt.Scanln(&che)
	if che == 1 || che == 3 {
		if che == 1 {
			var pos int
			fmt.Println("enter age : ")
			fmt.Scan(&pos)
			user := []User{}
			db.Where("Age = ?", pos).Find(&user)
			for i := 0; i < len(user); i++ {
				fmt.Print("User is : ", i+1, "\nname is =>", user[i].Name, "\nage is =>", user[i].Age, "\nid is =>", user[i].ID, "\nnid is =>", user[i].Nid, "\n\n")
			}
		} else {
			var pos int
			fmt.Println("enter id : ")
			fmt.Scan(&pos)
			user := User{}
			db.Where("ID = ?", pos).First(&user)
			fmt.Println("user is => ", user)
		}
	}
	if che == 2 || che == 4 {
		if che == 2 {
			fmt.Println("enter name : ")
			var posS string
			fmt.Scan(&posS)
			user := []User{}
			db.Where("Name = ?", posS).Find(&user)
			for i := 0; i < len(user); i++ {
				fmt.Print("User is : ", i+1, "\nname is =>", user[i].Name, "\nage is =>", user[i].Age, "\nid is =>", user[i].ID, "\nnid is =>", user[i].Nid, "\n\n")
			}
		} 
		if che == 4 {
			fmt.Println("enter Nid : ")
			var posS string
			fmt.Scan(&posS)
			user := User{}
			db.Where("Nid = ?", posS).First(&user)
			fmt.Println("user is => ", user)
		}
	}
}

func searchBook(db *gorm.DB) {
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

func loanBook(db *gorm.DB) {
	fmt.Print("enter id of book and user : ")
	var bookid, userid int
	fmt.Scan(&bookid, &userid)
	user := User{}
	book := Book{}
	db.Where("id = ?", userid).First(&user)
	if user.ID == userid {
		db.Where("b_id = ?", bookid).First(&book)
		if book.BID == bookid {
			fmt.Println("wait a moment...")
			loan := Loan{
				UID:  userid,
				BID:  bookid,
				Time: time.Now(),
			}
			db.Create(&loan)
			fmt.Println("loan seccess")
		} else {
			fmt.Println("not have this book.")
		}
	} else {
		fmt.Println("not have this user.")
	}
}

func AllUser(db *gorm.DB) {
	user := []User{}
	var can, i int64
	db.Find(&user).Count(&can)
	for i = 0; i < can; i++ {
		fmt.Print("user number : ", i+1, " is : \n", "user id is : ", user[i].ID, "\n", "user name is : ", user[i].Name, "\n", "user age is : ", user[i].Age, "\n", "user NID is : ", user[i].Nid, "\n\n")
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

func showUserLoan(db *gorm.DB){
	type Resualt struct{
		Name string
		ID int
		Title string
	}
	var resualt []Resualt
	db.Table("loans").Select("users.name, users.id, books.title").
	Joins("join users on users.id = loans.uid").
	Joins("join books on books.b_id = loans.b_id").
	Scan(&resualt)
	for i := 0; i < len(resualt); i++ {
		fmt.Printf("User Name: %s | ID: %d | Books : %s\n", resualt[i].Name, resualt[i].ID, resualt[i].Title)
	}
}

func getUserInput() int {
	var pos int
	fmt.Print("Enter your choice: ")
	fmt.Scan(&pos)
	return pos
}

func changeName(db *gorm.DB){
	user := User{}
	var id int
	var name string
	fmt.Println("enter your id and after that your name : ")
	fmt.Scan(&id, &name)
	db.Model(&user).Where("id = ?", id).Update("name", name)
}


func changeAge(db *gorm.DB){
	user := User{}
	var id, age int
	fmt.Println("enter your id and after that your age : ")
	fmt.Scan(&id, &age)
	db.Model(&user).Where("id = ?", id).Update("age", age)
}


func changeNID(db *gorm.DB){
	user := User{}
	fmt.Println(`
for National id i need all of your in formation so : 
1 : enter id
2 : enter name 
3 : enter old National id 
4 : your age.`)
	var nid, name string
	var id, age int
	fmt.Scan(&id, &name, &nid, &age)
	db.Where("id = ? and age = ? and nid = ? and name = ?", id, age, nid, name).First(&user)
	if user.ID == id && user.Age == age && user.Nid == nid && user.Name == name {
		var newNID string
		fmt.Print("Enter your new National ID: ")
		fmt.Scan(&newNID)
		db.Model(&user).Where("id = ? and age = ? and nid = ? and name = ?", id, age, nid, name).Update("nid", newNID)
	}
}

func changeAllInformation(db *gorm.DB){
	user := User{}
	fmt.Print("enter your natraul id : ")
	var nid string
	fmt.Scan(&nid)
	db.Where("nid", nid).First(&user)
	if user.Nid == nid{
		fmt.Print("ok enter your name and age for changing =>")
		var Nname string
		fmt.Scan(&Nname)
		var Nage int
		fmt.Scan(&Nage)
		db.Model(&user).Where("nid = ?", nid).Update("name", Nname)
		db.Model(&user).Where("nid = ?", nid).Update("age", Nage)
	}
}

func updateUser(db *gorm.DB) {
	fmt.Println(`
1 : change name =>
2 : change age => 
3 : change nid =>
4 : change all information =>`)
	switch os := getUserInput(); os{
	case 1: 
		changeName(db)
	case 2:
		changeAge(db)
	case 3:
		changeNID(db)
	case 4:
		changeAllInformation(db)
	}
}


func deleteUser(db *gorm.DB){
	var pos int
	fmt.Print("enter your pos : ")
	fmt.Print(`
1 : name
2 : Age
3 : id
4 : NId`)
	fmt.Print("\n")
	fmt.Scan(&pos)
	switch pos {
	case 1:
		var name string 
		fmt.Scan(&name)
		db.Where("name = ?", name).Delete(User{})
	case 2:
	case 3:
	case 4:
	}
}

func Lib(db *gorm.DB) {
	for on := 1; on != 0; {
		printpos()
		switch pos := getUserInput(); pos {
		case 0:
			on = 0
		case 1:
			addBook(db)
		case 2:
			addUser(db)
		case 3:
			serchUser(db)
		case 4:
			searchBook(db)
		case 5:
			AddManager(db)
		case 7:
			loanBook(db)
		case 8:
			AllUser(db)
		case 9:
			BookLoans(db)
		case 10:
			showUserLoan(db)
		case 11:
			updateUser(db)	
		case 12:
			deleteUser(db)
		}
	}
}


func main() {
	db := connectDB()
	Migrate(db)
	Lib(db)
}