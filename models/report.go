package models

type Revenues struct {
	TotalRevenue float32               `json:"totalRevenue"`
	PaymentTypes []RevenuePaymentTypes `json:"paymentTypes"`
}

type RevenuePaymentTypes struct {
	PaymentId          uint32  `json:"paymentId"`
	PaymentName        string  `json:"name"`
	PaymentLogo        string  `json:"logo"`
	PaymentTotalAmount float32 `json:"totalAmount"`
}

type Solds struct {
	SoldsOrderProduct []SoldsProduct `json:"orderProducts"`
}

type SoldsProduct struct {
	ProductId           uint32  `json:"productId"`
	ProductName         string  `json:"name"`
	ProductTotalQtySold string  `json:"totalQty"`
	ProductTotalAmount  float32 `json:"totalAmount"`
}
