package mailer

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

func sendMail(to, subject, body string) error {
	emailWeb := os.Getenv("EMAIL_APP")
	emailAccess := os.Getenv("EMAIL_ACCESS")

	m := gomail.NewMessage()
	m.SetHeader("From", emailWeb)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, emailWeb, emailAccess)
	return dialer.DialAndSend(m)
}
func ValidateEmailDeveloper(email string, token string, account string) error {
	emailWeb := os.Getenv("EMAIL_APP")
	emailAcces := os.Getenv("EMAIL_ACCESS")

	m := gomail.NewMessage()
	m.SetHeader("From", emailWeb)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Invitacion para ser developer")

	linkRegistre := os.Getenv("HOST_WEB") + "/verific/developer?token=" + token

	bodyHTML := fmt.Sprintf(`
		<h2>¡Hola! Has sido invitado al sistema</h2>
		<p>Para crear tu cuenta de %s, haz clic en el siguiente enlace:</p>
		<br>
		<a href="%s" style="padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px;">
			Completar mi registro
		</a>
		<br><br>
		<p>Este enlace expirará en 2 horas.</p>
	`, strings.ToLower(account), linkRegistre)

	m.SetBody("text/html", bodyHTML)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, emailWeb, emailAcces)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
func ValidateJob(email string, tenant_name, token, department, position string) error {
	linkRegistre := os.Getenv("HOST_WEB") + "/verific/developer?token=" + token

	body := fmt.Sprintf(`<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Invitación</title>
</head>
<body style="margin:0;padding:0;background:#f4f7fb;font-family:Arial,Helvetica,sans-serif;color:#333;">
  <table role="presentation" width="100%%" style="max-width:680px;margin:30px auto;background:#ffffff;border-radius:12px;box-shadow:0 6px 18px rgba(20,30,50,0.08);overflow:hidden;">
    <tr>
      <td style="padding:28px 32px 8px 32px;background:linear-gradient(90deg,#0ea5a4,#3b82f6);color:#fff;">
        <h1 style="margin:0;font-size:22px;line-height:1.2;">Bienvenido a <strong style="color:#fff;">%s</strong></h1>
        <p style="margin:8px 0 0 0;opacity:0.95;">Nos alegra que te unas a nuestro equipo.</p>
      </td>
    </tr>

    <tr>
      <td style="padding:22px 32px;">
        <p style="margin:0 0 14px 0;font-size:15px;color:#1f2937;">
          Has sido invitado a trabajar en el siguiente puesto. Por favor revisa los detalles y crea tu cuenta para unirte.
        </p>

        <div style="background:#f8fafc;border:1px solid #e6eef8;padding:14px;border-radius:8px;margin-bottom:18px;">
          <p style="margin:6px 0;font-size:14px;"><strong style="display:inline-block;width:110px;color:#0f172a;">Departamento:</strong> <span style="color:#0b5fff;">%s</span></p>
          <p style="margin:6px 0;font-size:14px;"><strong style="display:inline-block;width:110px;color:#0f172a;">Posición:</strong> <span style="color:#0b5fff;">%s</span></p>
        </div>

        <p style="margin:0 0 20px 0;font-size:14px;color:#475569;">
          Para completar tu registro y unirte al tenant, haz clic en el botón de abajo.
        </p>

        <p style="text-align:center;margin:0 0 18px 0;">
          <a href="%s" target="_blank" rel="noopener" style="display:inline-block;padding:12px 22px;background:linear-gradient(90deg,#06b6d4,#3b82f6);color:#fff;text-decoration:none;border-radius:8px;font-weight:600;box-shadow:0 6px 18px rgba(59,130,246,0.18);">
            Crear mi cuenta y unirme
          </a>
        </p>

        <hr style="border:none;border-top:1px solid #eef2f7;margin:18px 0;">

        <p style="font-size:12px;color:#94a3b8;margin:0;">
          Si el botón no funciona, copia y pega este enlace en tu navegador:
          <br><a href="%s" style="color:#2563eb;text-decoration:underline;">%s</a>
        </p>
      </td>
    </tr>

    <tr>
      <td style="padding:18px 32px;background:#fbfdff;text-align:center;font-size:12px;color:#9aa7bf;">
        <div>¿No esperabas esta invitación? Ignora este correo o contacta con tu administrador.</div>
      </td>
    </tr>
  </table>
</body>
</html>`, tenant_name, department, position, linkRegistre, linkRegistre, linkRegistre)

	return sendMail(email, "inicie secion para trabajar en "+tenant_name, body)
}

// Caso 2: El email YA existe (Solo debe aceptar o entrar) — sin rol
func SendInviteToOrganization(token, email, tenantName, department, position string) error {
	link := fmt.Sprintf("%s/login/verific/admin?token=%s", os.Getenv("HOST_WEB"), token)

	body := fmt.Sprintf(`<!doctype html>
<html lang="es">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Invitación a organización</title>
</head>
<body style="margin:0;padding:0;background:#f4f7fb;font-family:Arial,Helvetica,sans-serif;color:#333;">
  <table role="presentation" width="100%%" style="max-width:680px;margin:28px auto;background:#ffffff;border-radius:12px;box-shadow:0 6px 18px rgba(20,30,50,0.06);overflow:hidden;">
    <tr>
      <td style="padding:22px 28px;background:linear-gradient(90deg,#06b6d4,#3b82f6);color:#fff;">
        <h1 style="margin:0;font-size:20px;line-height:1.2;">Has sido añadido a <strong style="color:#fff;">%s</strong></h1>
        <p style="margin:8px 0 0 0;opacity:0.95;">Bienvenido de nuevo — ya tienes cuenta con nosotros.</p>
      </td>
    </tr>

    <tr>
      <td style="padding:20px 28px;">
        <p style="margin:0 0 14px 0;font-size:15px;color:#1f2937;">
          Hola, hemos detectado que ya tienes una cuenta. Te han invitado a la organización con los siguientes detalles:
        </p>

        <div style="background:#f8fafc;border:1px solid #e6eef8;padding:14px;border-radius:8px;margin-bottom:16px;">
          <p style="margin:6px 0;font-size:14px;">
            <strong style="display:inline-block;width:110px;color:#0f172a;">Organización:</strong>
            <span style="color:#0b5fff;">%s</span>
          </p>
          <p style="margin:6px 0;font-size:14px;">
            <strong style="display:inline-block;width:110px;color:#0f172a;">Departamento:</strong>
            <span style="color:#0b5fff;">%s</span>
          </p>
          <p style="margin:6px 0;font-size:14px;">
            <strong style="display:inline-block;width:110px;color:#0f172a;">Posición:</strong>
            <span style="color:#0b5fff;">%s</span>
          </p>
        </div>

        <p style="margin:0 0 18px 0;font-size:14px;color:#475569;">
          La próxima vez que inicies sesión verás esta empresa disponible en tu selector. Haz clic en el botón para ir a tu panel.
        </p>

        <p style="text-align:center;margin:0 0 18px 0;">
          <a href="%s" target="_blank" rel="noopener" style="display:inline-block;padding:12px 22px;background:linear-gradient(90deg,#06b6d4,#3b82f6);color:#fff;text-decoration:none;border-radius:8px;font-weight:600;box-shadow:0 6px 18px rgba(59,130,246,0.16);">
            Ir a mi panel
          </a>
        </p>

        <hr style="border:none;border-top:1px solid #eef2f7;margin:18px 0;">

        <p style="font-size:12px;color:#94a3b8;margin:0;">
          Si el botón no funciona, copia y pega este enlace en tu navegador:
          <br><a href="%s" style="color:#2563eb;text-decoration:underline;">%s</a>
        </p>
      </td>
    </tr>

    <tr>
      <td style="padding:16px 28px;background:#fbfdff;text-align:center;font-size:12px;color:#9aa7bf;">
        <div>¿No esperabas esta invitación? Ignora este correo o contacta con tu administrador.</div>
      </td>
    </tr>
  </table>
</body>
</html>`, tenantName, tenantName, department, position, link, link, link)

	return sendMail(email, "Fuiste añadido a "+tenantName, body)
}
