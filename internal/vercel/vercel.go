package vercel

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"

	"github.com/Implex-ltd/cleanhttp/cleanhttp"
	"github.com/Implex-ltd/fingerprint-client/fpclient"
	http "github.com/bogdanfinn/fhttp"
	"github.com/nu7hatch/gouuid"
)

func NewVercelClient() (*Client, error) {
	fp, err := fpclient.LoadFingerprint(&fpclient.LoadingConfig{
		FilePath: "../../assets/fingerprint.json",
	})

	if err != nil {
		return nil, err
	}

	c, err := cleanhttp.NewCleanHttpClient(&cleanhttp.Config{
		BrowserFp: fp,
	})

	if err != nil {
		return nil, err
	}

	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &Client{
		Http:        c,
		ChatUUID:    u.String(),
		Model:       "openai:gpt-3.5-turbo",
		MaxTokens:   500,
		Temperature: 1.7,
	}, nil
}

// this function is called into openai.jpeg
func (c *Client) Calculate(a float64) float64 {
	return a - math.Log10(a * math.Ln10)
}

func (c *Client) GetToken() (string, error) {
	response, err := c.Http.Do(cleanhttp.RequestOption{
		Method: "GET",
		Url:    "https://sdk.vercel.ai/openai.jpeg",
		Header: http.Header{
			`authority`:          {`sdk.vercel.ai`},
			`accept`:             {`*/*`},
			`accept-language`:    {`fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7`},
			`referer`:            {`https://sdk.vercel.ai/`},
			`sec-ch-ua`:          {`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
			`sec-ch-ua-mobile`:   {`?0`},
			`sec-ch-ua-platform`: {`"Windows"`},
			`sec-fetch-dest`:     {`empty`},
			`sec-fetch-mode`:     {`cors`},
			`sec-fetch-site`:     {`same-origin`},
			`user-agent`:         {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36`},

			http.HeaderOrderKey: {
				`authority`,
				`accept`,
				`accept-language`,
				`referer`,
				`sec-ch-ua`,
				`sec-ch-ua-mobile`,
				`sec-ch-ua-platform`,
				`sec-fetch-dest`,
				`sec-fetch-mode`,
				`sec-fetch-site`,
				`user-agent`,
			},
		},
	})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return "", err
	}

	var data B64Payload
	err = json.Unmarshal(decodedBytes, &data)
	if err != nil {
		return "", err
	}
	
	fmt.Println(data)
	result := c.Calculate(data.A)

	fmt.Println(fmt.Sprintf(`{"r":[%f,[],"mark"],"t":"%s"}`, result, data.T))

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`{"r":[%f,[],"mark"],"t":"%s"}`, result, data.T))), nil
}

func (c *Client) GetPrompt(prompt string) (string, error) {
	token, err := c.GetToken()
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(&SendPromptPayload{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		PlaygroundID:     c.ChatUUID,
		ChatIndex:        0,
		Model:            c.Model,
		Temperature:      c.Temperature,
		MaxTokens:        c.MaxTokens,
		TopK:             1,
		TopP:             1,
		FrequencyPenalty: 1,
		PresencePenalty:  1,
		StopSequences:    []string{},
	})

	if err != nil {
		return "", err
	}

	fmt.Println(string(payload), token)

	response, err := c.Http.Do(cleanhttp.RequestOption{
		Method: "POST",
		Url:    "https://sdk.vercel.ai/api/generate",
		Header: http.Header{
			`authority`:          {`sdk.vercel.ai`},
			`accept`:             {`*/*`},
			`accept-language`:    {`fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7`},
			`content-type`:       {`application/json`},
			`custom-encoding`:    {token},
			`origin`:             {`https://sdk.vercel.ai`},
			`referer`:            {`https://sdk.vercel.ai/`},
			`sec-ch-ua`:          {`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
			`sec-ch-ua-mobile`:   {`?0`},
			`sec-ch-ua-platform`: {`"Windows"`},
			`sec-fetch-dest`:     {`empty`},
			`sec-fetch-mode`:     {`cors`},
			`sec-fetch-site`:     {`same-origin`},
			`user-agent`:         {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36`},

			http.HeaderOrderKey: {
				`authority`,
				`accept`,
				`accept-language`,
				`content-type`,
				`custom-encoding`,
				`origin`,
				`referer`,
				`sec-ch-ua`,
				`sec-ch-ua-mobile`,
				`sec-ch-ua-platform`,
				`sec-fetch-dest`,
				`sec-fetch-mode`,
				`sec-fetch-site`,
				`user-agent`,
			},
		},
		Body: bytes.NewReader(payload),
	})

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
