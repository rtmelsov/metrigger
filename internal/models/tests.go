package models

type JSONReqType struct {
	T     string
	Name  string
	Delta int64
	Value float64
}

type JSONTest struct {
	Name       string
	Method     string
	ExpectBody string
	ExpectCode int
	Action     string
	Value      JSONReqType
}

type GetPingWebhook = []struct {
	Name       string
	Method     string
	ExpectCode int
}

type PostWebhookValue struct {
	T      string
	Name   string
	Number float64
}
type PostWebhook = []struct {
	Name       string
	Method     string
	ExpectCode int
	Value      PostWebhookValue
}

type GetWebhook = []struct {
	Name        string
	Method      string
	ContentType string
	ExpectCode  int
	Url         string
}
