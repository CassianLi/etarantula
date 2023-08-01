package service

import (
	"github.com/spf13/viper"
	"log"
	"strings"
	"tarantula-v2/utils"
)

// SendWaringEmail send waring email
func SendWaringEmail(salesChannel string, productNo string, weblink string, errs []string) error {
	body := viper.GetString("waring-email.body")

	body = strings.ReplaceAll(body, "{SALES_CHANNEL}", salesChannel)
	body = strings.ReplaceAll(body, "{PRODUCT_NO}", productNo)
	body = strings.ReplaceAll(body, "{WEBLINK}", weblink)
	errsStr := strings.Join(errs, "</br>")
	body = strings.ReplaceAll(body, "{ERRORS}", errsStr)

	emailTo := viper.GetStringSlice("waring-email.to")
	emailCc := viper.GetStringSlice("waring-email.cc")
	emailSubject := viper.GetString("waring-email.subject")

	log.Println("Send waring email...")
	log.Println("emailTo:", emailTo)
	log.Println("emailCc:", emailCc)
	log.Println("emailSubject:", emailSubject)
	log.Println("emailBody:", body)

	return utils.SendEmail(emailTo, emailCc, emailSubject, body)
}
