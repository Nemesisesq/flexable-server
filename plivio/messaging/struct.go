package messaging

type Response struct {
	Message Message
}
type Message struct {
	Src            string `xml:"src,attr"`
	Dst            string `xml:"dst,attr"`
	Type           string `xml:"type,attr"`
	CallbackUrl    string `xml:"callbackUrl,attr"`
	CallbackMethod string `xml:"callbackMethod,attr"`
	Value          string `xml:",innerxml"`
}
