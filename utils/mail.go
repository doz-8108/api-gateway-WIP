	package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendUserVerEmailReqBody struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	TemplateUuid      string            `json:"template_uuid"`
	TemplateVariables map[string]string `json:"template_variables"`
}

func (u *Utils) SendSignUpEmail(sender string, recipientName string,
	recipientEmail string,
	token string,
	redirectTo string,
) error {
	var reqBody SendUserVerEmailReqBody
	reqBody.From.Email = "no-reply@" + u.EnvVars.MAILTRAP_EMAIL_HOST
	reqBody.From.Name = u.EnvVars.MAILTRAP_EMAIL_HOST
	reqBody.To = append(reqBody.To, struct {
		Email string `json:"email"`
	}{Email: recipientEmail})
	reqBody.TemplateUuid = u.EnvVars.MAILTRAP_TEMPLATE_UUID
	reqBody.TemplateVariables = map[string]string{
		"name": recipientName,
		"url":  fmt.Sprintf("%s/users/verify?token=%s&redirect_to=%s", sender, token, redirectTo),
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", u.EnvVars.MAILTRAP_API_ENDPOINT, bytes.NewReader(payload))

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", u.EnvVars.MAILTRAP_API_TOKEN))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	// b, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(string(b))
	// }

	return nil
}
