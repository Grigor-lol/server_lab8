package main

import (
	"awesomeProject/API"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	api, err := API.NewAPI()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = api.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
