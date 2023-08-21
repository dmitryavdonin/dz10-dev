package model

type PaymentStatus struct {
	OrderId int    `json:"order_id,omitempty"`
	UserId  int    `json:"user_id,omitempty"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
}

type PaymentStatusMsg struct {
	Data PaymentStatus `json:"data"`
}
