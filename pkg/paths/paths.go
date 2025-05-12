package paths

import "strings"

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

func SplitFullName(fullName string) (string, string) {
	// fullName'i boşluklara göre ayır
	parts := strings.Fields(fullName)

	// Eğer birden fazla kelime varsa, ilk kısmı name ve son kısmı surname yapalım
	if len(parts) > 1 {
		name := strings.Join(parts[:len(parts)-1], " ") // İlk kelimeleri name olarak birleştir
		surname := parts[len(parts)-1]                  // Son kelimeyi surname olarak al
		return name, surname
	}

	// Eğer sadece bir kelime varsa, onu name olarak alalım, surname boş kalsın
	return parts[0], ""
}
