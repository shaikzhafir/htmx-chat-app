package models

type Message struct {
	User      string `json:"user"`
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
}

type HTMXMessage struct {
	Message string  `json:"message"`
	Headers Headers `json:"HEADERS"`
}

type Headers struct {
	HXRequest     string `json:"HX-Request"`
	HXTarget      string `json:"HX-Target"`
	HXTrigger     string `json:"HX-Trigger"`
	HXTriggerName string `json:"HX-Trigger-Name"`
	HXCurrentURL  string `json:"HX-Current-URL"`
}
