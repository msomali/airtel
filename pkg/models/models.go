package models

type AirtelAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

type AirtelAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type AirtelPushRequest struct {
	Reference  string `json:"reference"`
	Subscriber struct {
		Country  string `json:"country"`
		Currency string `json:"currency"`
		Msisdn   int64  `json:"msisdn"`
	} `json:"subscriber"`
	Transaction struct {
		Amount   int    `json:"amount"`
		Country  string `json:"country"`
		Currency string `json:"currency"`
		ID       string `json:"id"`
	} `json:"transaction"`
}

type AirtelPushResponse struct {
	Data struct {
		Transaction struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type AirtelRefundRequest struct {
	Transaction struct {
		AirtelMoneyID string `json:"airtel_money_id"`
	} `json:"transaction"`
}

type AirtelRefundResponse struct {
	Data struct {
		Transaction struct {
			AirtelMoneyID string `json:"airtel_money_id"`
			Status        string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type AirtelPushEnquiryResponse struct {
	Data struct {
		Transaction struct {
			AirtelMoneyID string `json:"airtel_money_id"`
			ID            string `json:"id"`
			Message       string `json:"message"`
			Status        string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type AirtelCallbackRequest struct {
	Transaction struct {
		ID            string `json:"id"`
		Message       string `json:"message"`
		StatusCode    string `json:"status_code"`
		AirtelMoneyID string `json:"airtel_money_id"`
	} `json:"transaction"`
	Hash string `json:"hash"`
}

type AirtelUserEnquiryResponse struct {
	Data struct {
		FirstName    string `json:"first_name"`
		Grade        string `json:"grade"`
		IsBarred     bool   `json:"is_barred"`
		IsPinSet     bool   `json:"is_pin_set"`
		LastName     string `json:"last_name"`
		Msisdn       int    `json:"msisdn"`
		RegStatus    string `json:"reg_status"`
		RegisteredID string `json:"registered_id"`
		Registration struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"registration"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type AirtelBalanceEnquiry struct {
	Data struct {
		Balance       string `json:"balance"`
		Currency      string `json:"currency"`
		AccountStatus string `json:"account_status"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type AirtelDisburseRequest struct {
	Payee struct {
		Msisdn string `json:"msisdn"`
	} `json:"payee"`
	Reference   string `json:"reference"`
	Pin         string `json:"pin"`
	Transaction struct {
		Amount int    `json:"amount"`
		ID     string `json:"id"`
	} `json:"transaction"`
}

type AirtelDisburseEnquiry struct {
	Data struct {
		Transaction struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"transaction"`
	} `json:"data"`
	Status struct {
		Code       string `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

type ListTransactionsResponse struct {
	Data struct {
		Count        int `json:"count"`
		Transactions []struct {
			Charges struct {
				Service int `json:"service"`
			} `json:"charges"`
			Payee struct {
				Currency string `json:"currency"`
				Msisdn   int    `json:"msisdn"`
				Name     string `json:"name"`
			} `json:"payee"`
			Payer struct {
				Currency string `json:"currency"`
				Msisdn   int    `json:"msisdn"`
				Name     string `json:"name"`
			} `json:"payer"`
			Service struct {
				Type string `json:"type"`
			} `json:"service"`
			Transaction struct {
				AirtelMoneyID   string `json:"airtel_money_id"`
				Amount          int    `json:"amount"`
				CreatedAt       int    `json:"created_at"`
				ID              int64  `json:"id"`
				ReferenceNumber string `json:"reference_number"`
				Status          string `json:"status"`
			} `json:"transaction"`
		} `json:"transactions"`
	} `json:"data"`
	Status struct {
		Code       int    `json:"code"`
		Message    string `json:"message"`
		ResultCode string `json:"result_code"`
		Success    bool   `json:"success"`
	} `json:"status"`
}

