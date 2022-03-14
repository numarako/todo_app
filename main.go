package main

import (
	"fmt"
	"todo_app/app/controllers"
	"todo_app/app/models"
)

func main() {
	i := 9
	fmt.Println(i)
	fmt.Println(models.Db)

	controllers.StartMainServer()
}
