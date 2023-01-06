package model

type Settings struct {
	AwsAccessKeyId           string `json:"aws_access_key_id"`
	AwsSecretAccessKey       string `json:"aws_secret_access_key"`
	AwsRegion                string `json:"aws_region"`
	StripePubKey             string `json:"stripe_pub_key"`
	StripePrivateKey         string `json:"stripe_private_key"`
	StripeTestPubKey         string `json:"stripe_test_pub_key"`
	StripeTestPrivateKey     string `json:"stripe_test_private_key"`
	StripeMode               string `json:"stripe_mode"`
	SmtpHost                 string `json:"smtp_host"`
	SmtpPort                 string `json:"smtp_port"`
	SmtpUser                 string `json:"smtp_user"`
	SmtpPassword             string `json:"smtp_password"`
	SmtpTls                  string `json:"smtp_tls"`
	GoogleServiceAccountJson string `json:"google_service_account_json"`
}
