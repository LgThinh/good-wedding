package main

import (
	"good-wedding/conf"
	"good-wedding/pkg/router"
)

// @securityDefinitions.apikey Authorization
// @in                         header
// @name                       Authorization

// @securityDefinitions.apikey User ID
// @in                         header
// @name                       x-user-id
func main() {
	conf.LoadConfig()
	router.NewRoute()
}
