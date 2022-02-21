package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Accrual struct {
	Address string
}

type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

var ErrNotFound = errors.New("order not found")

type ErrTooManyRequests struct {
	Err        error
	Message    string
	RetryAfter *time.Time
}

func (e ErrTooManyRequests) Unwrap() error {
	return e.Err
}

func (e ErrTooManyRequests) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "too many requests"
}

func (a Accrual) GetOrderStatus(ctx context.Context, number string) (*OrderStatus, error) {
	url := fmt.Sprintf("%s/api/orders/%s", a.Address, number)
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	response, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		retryValue := response.Header.Get("retry-after")
		retryAfter := a.parseRetryAfterValue(retryValue)

		return nil, ErrTooManyRequests{
			Message:    string(body),
			RetryAfter: retryAfter,
		}
	} else if response.StatusCode == http.StatusNoContent || response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	orderStatus := &OrderStatus{}
	err = json.Unmarshal(body, orderStatus)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

func (a Accrual) parseRetryAfterValue(retryValue string) *time.Time {
	if retryValue == "" {
		return nil
	}

	if sec, err := strconv.Atoi(retryValue); err == nil {
		// this is seconds
		after := time.Now().Add(time.Duration(sec) * time.Second)
		return &after
	}
	// this is time
	after, err := time.Parse(time.RFC3339, retryValue)
	if err != nil {
		return nil
	}

	return &after
}
