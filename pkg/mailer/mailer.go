package mailer

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func ValidateEmailDeveloper(email string, token string, account string) error {
	err:=godotenv.Load()
	if err!=nil {
		return fmt.Errorf("Error load the archive .env: %v",err)
	}
	emailWeb:=os.Getenv("EMAIL_APP")
	emailAcces:=os.Getenv("EMAIL_ACCESS")

	m := gomail.NewMessage()
	m.SetHeader("From",emailWeb)
	m.SetHeader("To",email)
	m.SetHeader("Subject","Invitacion para ser developer")
	
	linkRegistre:=os.Getenv("HOST_WEB")+"/verific/developer?token="+token

	bodyHTML := fmt.Sprintf(`
		<h2>¡Hola! Has sido invitado al sistema</h2>
		<p>Para crear tu cuenta de %s, haz clic en el siguiente enlace:</p>
		<br>
		<a href="%s" style="padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px;">
			Completar mi registro
		</a>
		<br><br>
		<p>Este enlace expirará en 2 horas.</p>
	`,strings.ToLower(account), linkRegistre)

	m.SetBody("text/html",bodyHTML)

	dialer:=gomail.NewDialer("smtp.gmail.com",587,emailWeb,emailAcces)

	if err:=dialer.DialAndSend(m);err!=nil {
		return err
	}
	return nil
}