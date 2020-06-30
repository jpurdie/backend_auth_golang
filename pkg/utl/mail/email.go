package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jpurdie/authapi"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"os"
	"text/template"
	"time"
)

func buildEmailClient() *mailgun.MailgunImpl {
	return mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"))
}

func SendInvitationEmail(i *authapi.Invitation) error {
	mg := buildEmailClient()

	webURL := os.Getenv("APP_WEB")

	data := struct {
		Token  string
		WebURL string
	}{
		Token:  i.TokenStr,
		WebURL: webURL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	t := template.New("invitation.html")

	var err error
	t, err = t.ParseFiles("pkg\\utl\\mail\\templates\\invitation.html")
	if err != nil {
		log.Println(err)

		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		log.Println(err)
		return err
	}

	result := tpl.String()

	sender := os.Getenv("SENDER_EMAIL")
	subject := "You've been invited to Vitae!"
	recipient := i.Email

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, result, recipient)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	message.SetHtml(result)

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Println(err)

		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return err
}
