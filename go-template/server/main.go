package main

import (
	"good-template-go/conf"
	"good-template-go/pkg/service"
)

// @securityDefinitions.apikey Authorization
// @in                         header
// @name                       Authorization

// @securityDefinitions.apikey User ID
// @in                         header
// @name                       x-user-id
func main() {
	conf.LoadConfig()
	service.NewService()
}
