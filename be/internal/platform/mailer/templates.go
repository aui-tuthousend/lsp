package mailer

import "fmt"

func RenderAssessorAssigned(name, link string) EmailMessage {
	text := fmt.Sprintf("Halo %s, Anda mendapat penugasan asesmen baru. Buka: %s", name, link)
	return EmailMessage{Subject: "Penugasan Asesor Baru", TextBody: text, HTMLBody: fmt.Sprintf("<p>%s</p><p><a href=\"%s\">Buka penugasan</a></p>", text, link)}
}

func RenderAPL02Submitted(name, link string) EmailMessage {
	text := fmt.Sprintf("Halo %s, asesi telah submit APL-02. Tinjau di: %s", name, link)
	return EmailMessage{Subject: "APL-02 Siap Ditinjau", TextBody: text, HTMLBody: fmt.Sprintf("<p>%s</p><p><a href=\"%s\">Tinjau APL-02</a></p>", text, link)}
}
