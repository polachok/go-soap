package soap

import "encoding/xml"
import "net/http"
import "bytes"
import "io"

type Envelope struct {
	XMLName       xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Xsi           string   `xml:"xmlns:xsi,attr"`
	Soapenc       string   `xml:"xmlns:soapenc,attr"`
	Xsd           string   `xml:"xmlns:xsd,attr"`
	EncodingStyle string   `xml:"soap:encodingStyle,attr"`
	Soap          string   `xml:"xmlns:soap,attr"`
	Body          Body
}

type Body struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Data    string   `xml:",innerxml"`
}

func NewEnvelope(data interface{}) Envelope {
	msg, err := xml.Marshal(data)
	if err != nil {
		panic(err)
	}
	return Envelope{
		Xsi:           "http://www.w3.org/2001/XMLSchema-instance",
		Soapenc:       "http://schemas.xmlsoap.org/soap/encoding/",
		Xsd:           "http://www.w3.org/2001/XMLSchema",
		EncodingStyle: "http://schemas.xmlsoap.org/soap/encoding/",
		Soap:          "http://schemas.xmlsoap.org/soap/envelope/",
		Body:          Body{Data: string(msg)},
	}
}

func WriteEnvelope(env Envelope, writer io.Writer) error {
	msg, err := xml.Marshal(env)
	if err != nil {
		return err
	}

	writer.Write(msg)
	return nil
}

func SendEnvelope(env Envelope, url string) (error, *http.Response) {
	buf := new(bytes.Buffer)

	buf.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	err := WriteEnvelope(env, buf)
	if err != nil {
		return err, nil
	}

	r, err := http.Post(url, "application/soap+xml", buf)
	if err != nil {
		return err, nil
	} else {
		return nil, r
	}
}

func ReadEnvelope(reply interface{}, reader io.Reader) error {
	var env Envelope

	dec := xml.NewDecoder(reader)
	err := dec.Decode(&env)
	if err != nil {
		return err
	}
	return xml.Unmarshal([]byte(env.Body.Data), reply)
}
