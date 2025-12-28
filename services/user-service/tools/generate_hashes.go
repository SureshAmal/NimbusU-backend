package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := map[string]string{
		"Admin":   "Admin@123",
		"Faculty": "Faculty@123",
		"Student": "Student@123",
		"Staff":   "Staff@123",
	}

	fmt.Println("Generating bcrypt password hashes for NimbusU seed data...")
	fmt.Println("========================================================")
	fmt.Println()

	for role, password := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash password for %s: %v", role, err)
		}
		fmt.Printf("%s (%s):\n%s\n\n", role, password, hash)
	}

	fmt.Println("========================================================")
	fmt.Println("Copy these hashes and update the seed.sql file")
	fmt.Println("Replace the placeholder password_hash values with the generated hashes")
}
