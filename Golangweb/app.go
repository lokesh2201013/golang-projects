package main

import (
	"fmt"
	"log"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	fmt.Println("hello world")
	fmt.Printf("Current user: %s\n", usr.Username)
	fmt.Printf("User ID: %s\n", usr.Uid)
	fmt.Printf("Group ID: %s\n", usr.Gid)
	fmt.Printf("Home Directory: %s\n", usr.HomeDir)
}
