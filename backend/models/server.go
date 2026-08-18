package models

type SMTPData struct {
	Username string
	Password string
	Host     string
	Addr     string
	From     string
	BaseURL  string
}

func NewSMTPData(username, password, host, addr, from, baseURL string) SMTPData {
	return SMTPData{
		Username: username,
		Password: password,
		Host:     host,
		Addr:     addr,
		From:     from,
		BaseURL:  baseURL,
	}
}

type ServerData struct {
	URI      string
	BaseURL  string
	SMTPData SMTPData
}

func NewServerData(URI string, SMTPData SMTPData) ServerData {
	return ServerData{
		URI:      URI,
		SMTPData: SMTPData,
	}
}
