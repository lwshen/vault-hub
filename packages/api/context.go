package api

// ClientInfo captures request metadata used for auditing.
type ClientInfo struct {
	IP        string
	UserAgent string
}

// ExtractClientInfo returns client metadata based on common forwarding headers.
func ExtractClientInfo(headerGetter func(string) string, remoteIP string) ClientInfo {
	ip := headerGetter("X-Forwarded-For")
	if ip == "" {
		ip = headerGetter("X-Real-IP")
	}
	if ip == "" {
		ip = remoteIP
	}

	return ClientInfo{
		IP:        ip,
		UserAgent: headerGetter("User-Agent"),
	}
}
