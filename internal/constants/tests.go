package constants

import "github.com/rtmelsov/metrigger/internal/models"

var JSONTests = []models.JSONTest{
	{
		Name:       "1",
		Action:     "update",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":3242, "id":"jsonTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:     "counter",
			Name:  "jsonTest",
			Delta: 3242,
		},
	},
	{
		Name:       "2",
		Action:     "update",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":6484, "id":"jsonTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:     "counter",
			Name:  "jsonTest",
			Delta: 3242,
		},
	},
	{
		Name:       "3",
		Action:     "update",
		ExpectBody: `{"id":"jsonTest", "type":"gauge", "value":32.42}`,
		Method:     "POST",
		ExpectCode: 200,
		Value: models.JSONReqType{
			T:     "gauge",
			Name:  "jsonTest",
			Value: 32.42,
		},
	},
	{
		Name:       "4",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":6484, "id":"jsonTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:    "counter",
			Name: "jsonTest",
		},
	},
	{
		Name:       "5",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"id":"jsonTest", "type":"gauge", "value":32.42}`,
		Value: models.JSONReqType{
			T:    "gauge",
			Name: "jsonTest",
		},
	},
	{
		Name:       "6",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 404,
		ExpectBody: "",
		Value: models.JSONReqType{
			T:    "gauge",
			Name: "unknown",
		},
	},
}

var GzipTests = []models.JSONTest{
	{
		Name:       "1",
		Action:     "update",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":3242, "id":"gzipTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:     "counter",
			Name:  "gzipTest",
			Delta: 3242,
		},
	},
	{
		Name:       "2",
		Action:     "update",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":6484, "id":"gzipTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:     "counter",
			Name:  "gzipTest",
			Delta: 3242,
		},
	},
	{
		Name:       "3",
		Action:     "update",
		ExpectBody: `{"id":"gzipTest", "type":"gauge", "value":32.42}`,
		Method:     "POST",
		ExpectCode: 200,
		Value: models.JSONReqType{
			T:     "gauge",
			Name:  "gzipTest",
			Value: 32.42,
		},
	},
	{
		Name:       "4",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"delta":6484, "id":"gzipTest", "type":"counter"}`,
		Value: models.JSONReqType{
			T:    "counter",
			Name: "gzipTest",
		},
	},
	{
		Name:       "5",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 200,
		ExpectBody: `{"id":"gzipTest", "type":"gauge", "value":32.42}`,
		Value: models.JSONReqType{
			T:    "gauge",
			Name: "gzipTest",
		},
	},
	{
		Name:       "6",
		Action:     "value",
		Method:     "POST",
		ExpectCode: 404,
		ExpectBody: "",
		Value: models.JSONReqType{
			T:    "gauge",
			Name: "unknown",
		},
	},
}

var GetPingWebhook = models.GetPingWebhook{
	{
		Name:       "ping",
		Method:     "GET",
		ExpectCode: 200,
	},
}

var PostWebhook = models.PostWebhook{{
	Name:       "1",
	Method:     "POST",
	ExpectCode: 200,
	Value: models.PostWebhookValue{
		T:      "counter",
		Name:   "fdsafd",
		Number: 3242,
	},
},
	{
		Name:       "2",
		Method:     "POST",
		ExpectCode: 200,
		Value: models.PostWebhookValue{
			T:      "gauge",
			Name:   "fdsafd",
			Number: 3242,
		},
	},
}

var GetWebhook = models.GetWebhook{
	{
		Name:        "1",
		Method:      "GET",
		ContentType: "application/json",
		ExpectCode:  200,
		Url:         "/value/counter/fdsafd",
	},
	{
		Name:        "2",
		Method:      "GET",
		ContentType: "application/json",
		ExpectCode:  200,
		Url:         "/value/gauge/fdsafd",
	},
	{
		Name:        "3",
		Method:      "GET",
		ContentType: "text/html",
		ExpectCode:  200,
		Url:         "/",
	},
}
