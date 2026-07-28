package main

import "github.com/viettrung2103/bookmark-management/internal/infrastructure"

// @title Bookmark API
// @version 4.0.0
// @description API for bookmark management
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @BasePath /
// @in header
// @name Authorization
func main() {

	infrastructure.StartApp()
}
