package listennotes

import (
	"fmt"
)

// Known errors
var (
	ErrBadRequest          = fmt.Errorf("something wrong on your end (client side errors), e.g., missing required parameters")
	ErrUnauthorized        = fmt.Errorf("wrong api key or your account is suspended")
	ErrNotFound            = fmt.Errorf("endpoint does not exist, or podcast / episode does not exist")
	ErrTooManyRequests     = fmt.Errorf("for FREE plan, exceeding the quota limit; or for all plans, sending too many requests too fast and exceeding the rate limit - https://www.listennotes.com/api/faq/#faq17")
	ErrInternalServerError = fmt.Errorf("something wrong on our end (unexpected server errors)")
)

var errMap = map[int]error{
	200: nil,
	400: ErrBadRequest,
	401: ErrUnauthorized,
	404: ErrNotFound,
	429: ErrTooManyRequests,
	500: ErrInternalServerError,
}
