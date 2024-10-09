package utils

import (
	"fmt"
	"time"

	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/config"
	"github.com/Inter-IIT-Prepathon-TheSloths/backend/internal/services"
	"github.com/go-mail/mail"
	"github.com/labstack/echo/v4"
)

func SendVerification(c echo.Context, sv *services.UserService, email, subject, heading, info1, time_duration, regenerate_link string, t time.Duration, extras map[string]string) error {

	code, err := GenerateOTP()
	if err != nil {
		return err
	}

	if err = sv.UpdateVerification(c.Request().Context(), email, code, time.Now().Add(t), extras); err != nil {
		return err
	}

	from := config.FromEmail
	password := config.SmtpAppPass

	logo_url := "https://iili.io/dQyCGGj.png"

	body := fmt.Sprintf(`
		<html lang="en">

		<body>
			<table align="center" width="100%%" cellpadding="0" cellspacing="0">
				<tr>
					<td align="center">
						<img src="%s" width="70px" alt="Sloth">
					</td>
				</tr>
				<tr>
					<td align="center">
						<h2>%s</h2>
					</td>
				</tr>
				<tr>
					<td align="center">
						<table align="center" width="500" height="300" cellpadding="0" cellspacing="0"
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
									<span
										style="text-align: center; background-color: rgb(201, 201, 201); border-radius: 2px; height: 40px; padding: 10px; margin: 10px; font-weight: bold; margin: 10px">
										%s</span>
								</td>
							</tr>
							<tr>
								<td align="center">
									<p>If you don't use this link within %s, it will expire. Click <a target="_blank" href="%s">here</a> to get a new
										account activation link.</p>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>

		</html>
	`, logo_url, heading, info1, code, time_duration, regenerate_link)

	m := mail.NewMessage()

	m.SetHeader("From", from)

	m.SetHeader("To", email)

	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	d := mail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
