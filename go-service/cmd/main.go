package main

import (

	router "github.com/SLANGERES/go-service/internal/Routers"

)

func main(){
	router:=router.Router()

	router.Run(":9090")
	
}