package models

type SMTPData struct {
	Username string
	Password string
	Host string
	Addr string
	From string
}

func NewSMTPData(username,password,host,addr,from string) SMTPData {
	return SMTPData {
		Username: username,
		Password: password,
		Host: host,
		Addr: addr,
		From: from,
	}
}

type ServerData struct {
	URI string
	BaseURL string
	SMTPData SMTPData
}

func NewServerData(URI,baseURL string,SMTPData SMTPData) ServerData {
	return ServerData {
		URI:URI,
		BaseURL:baseURL,
		SMTPData: SMTPData,
	}
}

