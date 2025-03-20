package handlers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
)

func sendEmail(name, email, message string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	// Create a new message
	msg := gomail.NewMessage()

	// Set email headers
	msg.SetHeader("From", email)
	msg.SetHeader("To", "aramalho.1991@gmail.com")
	msg.SetHeader("Subject", "Vi seu site e gostaria de entrar em contato")

	// Set email body
	msg.SetBody("text/plain", message)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer("smtp.gmail.com", 587, user, password)

	// Send the email
	if err := dialer.DialAndSend(msg); err != nil {
		log.Println("Error sending email:", err)
		return err
	}
	log.Println("Email sent successfully!")
	return nil
}
