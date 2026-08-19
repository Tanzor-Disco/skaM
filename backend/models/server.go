package models

// SMTPData represents data required for go SMTP package
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

// ServerData represents a set of data passed to a new instance of server struct
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
