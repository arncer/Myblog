package main

import (
	"exchangeapp/config"
	"exchangeapp/router"
	"fmt"
)

func main() {
	config.InitConfig()
	r := router.SetUpRouter()
	port := config.AppConfig.App.Port

	if port == "" {
		port = "8080"
	}
	

	err := r.Run(config.AppConfig.App.Port) // listen and serve on 0.0.0.0:8080
	if err != nil {
		fmt.Println(err.Error())
	}
	

}
