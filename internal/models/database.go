package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" 
)




var DB *sql.DB

func InitializeDB() error{
	var err error

	DB,err=sql.Open("sqlite3","./forum.db")
	if err != nil{
		return fmt.Errorf("failed to open database: %v",err)
	}

	if err := createTables();err !=nil{
		return fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error{
	for _,query:= range Tables{
		_,err:=DB.Exec(query)
		if err !=nil{
			return fmt.Errorf("failed to execute query: %v",err)
		}
	}
	return nil
}