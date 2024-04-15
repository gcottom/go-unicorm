package main

import "fmt"

func main() {
	c := CreateClient()
	p := Person{ID: "1", FirstName: "Gage", LastName: "Clockjaw", Age: 29, Children: 0}
	c.PersonTable.Save(p)
	fmt.Println(c.PersonTable.GetAll())

}
