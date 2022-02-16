package requests

type BalanceWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
