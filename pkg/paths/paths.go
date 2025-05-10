package paths

func CreateURI(http bool, id, domain string) string {
	return CreateURL(http, domain) + "metadata/" + id
}

func CreateURL(http bool, domain string) string {
	var uri string
	if http {
		uri += "https://" + domain + "/"
	} else {
		uri += "http://" + domain + "/"
	}

	return uri
}
