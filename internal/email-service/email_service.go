package email

import "fmt"

// for now
type EmailService struct {
}

func (es *EmailService) SendEmail(content string) error {
	fmt.Printf("== Sening email: %s\n", content)
	return nil
}
