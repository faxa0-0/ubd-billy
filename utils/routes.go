package utils

import (
	"fmt"
	"strings"
)

type route struct {
	method  string
	url     string
	handler string
}
type groups struct {
	name   string
	routes []route
}

func LogRoutes() {
	base := "http://localhost:8080/api"
	routes := []groups{
		{
			name: "Auth",
			routes: []route{
				{method: "POST", url: "/login", handler: "LoginHandler"},
				{method: "POST", url: "/logout", handler: "LogoutHandler"},
				{method: "POST", url: "/refresh", handler: "RefreshHandler"}},
		},
		{
			name: "User",
			routes: []route{
				{method: "POST", url: "/users", handler: "CreateUserHandler"},
				{method: "GET", url: "/users", handler: "GetUsersHandler"},
				{method: "GET", url: "/users/{id}", handler: "GetUserByIDHandler"}},
		},
		{
			name: "Usage",
			routes: []route{
				{method: "GET", url: "/usage/{id}", handler: "GetUsageHandler"}},
		},
	}

	fmt.Println("\033[34m" + "Available routes" + "\033[0m")
	for _, group := range routes {
		fmt.Println(strings.Repeat("-", 88))
		fmt.Printf("\033[35m%45s\033[0m\n", group.name)
		fmt.Println(strings.Repeat("-", 88))
		for _, route := range group.routes {
			fmt.Printf("|%-10s | %-50s | %-20s|\n", route.method, base+route.url, route.handler)
		}
	}
	fmt.Println(strings.Repeat("-", 88))
}
