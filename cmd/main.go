package main

import (
	"fmt"
)

func main() {
	c := CreateClient()
	fmt.Println(c.PersonTable.GetByFirstName("Lori"))
	//fmt.Println(c.PersonTable.GetAll())

}
