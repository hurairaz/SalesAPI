package main

import (
	"SalesAPI/config"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type Author struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Book struct {
	gorm.Model
	Title    string
	AuthorID *uint
}

func main() {
	fmt.Println("Hello, World")
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err.Error())
	}
	fmt.Println(db)

	//err = db.Migrator().DropTable(&Author{}, &Book{})

	//err = db.AutoMigrate(&Author{}, &Book{})

	//var author = Author{Name: "James"}
	//err = db.Create(&author).Error
	//
	//var book = Book{Title: "The Great Gatsby", AuthorID: &author.ID}
	//err = db.Create(&book).Error

	//err = db.Delete(&Author{}, 1).Error
	//fmt.Println(err)

}
