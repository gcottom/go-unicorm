package main

import (
	"database/sql"
	"fmt"

	"github.com/gcottom/go-unicorm"

	_ "modernc.org/sqlite"
)

type Client struct {
	PersonTable *PersonTable
	DBClient    *sql.DB
}

func CreateClient() *Client {
	db, err := sql.Open("sqlite", "db.db")
	if err != nil {
		panic(err)
	}
	p := new(PersonTable)
	p.Table = new(unicorm.Table[Person])
	unicorm.InitTable(p, Person{}, db)

	return &Client{DBClient: db, PersonTable: p}
}

type PersonTable struct {
	*unicorm.Table[Person]
}

type Person struct {
	ID        string
	FirstName string
	LastName  string
	Age       int
	Test      string
	Children  int
}

func (p *PersonTable) GetAll() []Person {
	r, e := p.AutoGenerate()
	if e != nil {
		panic(e)
	}
	return r
}
func (p *PersonTable) GetByID(id string) Person {
	r, e := p.AutoGenerate(id)
	if e != nil {
		panic(e)
	}
	return r[0]
}
func (p *PersonTable) GetByFirstName(name string) Person {
	r, e := p.AutoGenerate(name)
	if e != nil {
		panic(e)
	}
	return r[0]
}

func (p *PersonTable) Save(person Person) error {
	_, err := p.AutoGenerate(person)
	return err
}

func (p Person) String() string {
	return fmt.Sprintf("ID:%s, Name:%s %s, Age:%d, Children:%d", p.ID, p.FirstName, p.LastName, p.Age, p.Children)
}
