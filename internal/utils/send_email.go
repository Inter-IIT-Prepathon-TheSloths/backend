package utils

import (
	"fmt"
	"os"

	"github.com/go-mail/mail"
)

func SendEmail(to []string, subject, heading, info1, link, time_duration, regenerate_link string) error {
	from := os.Getenv("FROM_EMAIL")
	password := os.Getenv("SMTP_APP_PASS")

	logo_url := "https://iili.io/dQyCGGj.png"
	// heading := "Activate your account"
	// info1 := "To activate your account, please click the button below and follow the instructions provided."
	// link := "random_link"
	// time_duration := "1 day"
	// regenerate_link := "random"

	body := fmt.Sprintf(`
		<html lang="en">

		<body>
			<table align="center" width="100%%" cellpadding="0" cellspacing="0">
				<tr>
					<td align="center">
						<img src=%s width="70px" alt="Sloth">
					</td>
				</tr>
				<tr>
					<td align="center">
						<h2>%s</h2>
					</td>
				</tr>
				<tr>
					<td align="center">
						<table align="center" width="500" cellpadding="0" cellspacing="0"
							style="border-radius: 10px; border: 1px solid rgb(195, 193, 193); padding: 30px;">
							<tr>
								<td align="center">
									<p>
										%s
									</p>
								</td>
							</tr>
							<tr>
								<td align="center">
									<a href="%s" target="_blank" style="text-decoration: none;">
										<button
											style="text-align: center; background-color: blue; border-radius: 15px; color: white; height: 40px; padding-left: 15px; padding-right: 15px; font-weight: bold; cursor: pointer; margin: 10px">Activate
											account</button>
									</a>
								</td>
							</tr>
							<tr>
								<td align="center">
									<p>
									If you don't use this link within %s, it will expire. Click <a target="_blank"
											href="%s">here</a> to get a new account activation link.</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>

		</html>
	`, logo_url, heading, info1, link, time_duration, regenerate_link)

	m := mail.NewMessage()

	m.SetHeader("From", from)

	m.SetHeader("To", to[0])

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := mail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
