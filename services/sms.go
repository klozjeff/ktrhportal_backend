package services

import (
	"bytes"
	"encoding/json"
	"io"
	"ktrhportal/utilities"
	"net/http"
)

type AdvantaResponse struct {
	Responses []interface{} `json:"responses"`
}

func SendSMS(phone_no string, message string) (string, error) {
	var respBulkSMS AdvantaResponse
	baseURL := utilities.GoDotEnvVariable("ADVANTA_BASE_URL")
	postBody, _ := json.Marshal(map[string]string{
		"apikey":    utilities.GoDotEnvVariable("ADVANTA_API_KEY"),
		"partnerID": utilities.GoDotEnvVariable("ADVANTA_PARTNER_ID"),
		"message":   message,
		"shortcode": utilities.GoDotEnvVariable("DEFAULT_SMS_SENDER_ID"),
		"mobile":    phone_no,
	})
	responseBody := bytes.NewBuffer(postBody)

	client := http.Client{}
	req, reqErr := http.NewRequest("POST", baseURL, responseBody)

	if reqErr != nil {

		return "", reqErr
	}
	req.Header = http.Header{
		"Content-Type":     []string{"application/json"},
		"X-Requested-With": []string{"XMLHttpRequest"},
	}
	resp, resErr := client.Do(req)

	if resErr != nil {
		return "", resErr
	}
	defer resp.Body.Close()

	body, bodyErr := io.ReadAll(resp.Body)

	if bodyErr != nil {
		return "", bodyErr
	}

	json.Unmarshal(body, &respBulkSMS)

	msgResp, _ := respBulkSMS.Responses[0].(map[string]interface{})

	messageid := msgResp["messageid"]

	return messageid.(string), nil
}
