package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
)

type Seed struct {
	DB *sql.DB
}

func Seeder(db *sql.DB) *Seed {
	return &Seed{
		DB: db,
	}
}

// user struct dùng để seed dữ liệu
type user struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Password   string
	Phone      string
	Status     string
}

func (m *Seed) SeedUsers() {
	users := []user{
		{
			FirstName:  "Thang",
			MiddleName: "Duc",
			LastName:   "Nguyen",
			Email:      "thangdz@gmail.com",
			Password:   randomAndWrite(16),
			Phone:      "0901234567",
			Status:     "active",
		},
	}

	for _, u := range users {
		// hash password
		password, err := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
		if err != nil {
			log.Fatalln(err)
		}

		// insert user
		_, err = m.DB.ExecContext(
			context.Background(),
			`INSERT INTO users 
				(first_name, middle_name, last_name, email, password_hash, phone, status, created_at, updated_at, verified_at) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
			u.FirstName,
			u.MiddleName,
			u.LastName,
			u.Email,
			password,
			u.Phone,
			u.Status,
			time.Now(),
			time.Now(),
			time.Now(),
		)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func writeToEnv(password string) {
	f, err := os.OpenFile(".env",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("\nADMIN_PASSWORD=" + password + "\n"); err != nil {
		log.Println(err)
	}
}

func randomAndWrite(n int) string {
	var chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()_+"

	ll := len(chars)
	b := make([]byte, n)
	_, _ = rand.Read(b)
	for i := 0; i < n; i++ {
		b[i] = chars[int(b[i])%ll]
	}

	str := string(b)
	fmt.Printf("Password is: %s\n", str)

	writeToEnv(str)

	return str
}
