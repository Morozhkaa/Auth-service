package models

import "fmt"

var (
	ErrBadRequest    = fmt.Errorf("required parameters are not filled") // 400
	ErrForbidden     = fmt.Errorf("forbidden: wrong password")          // 403
	ErrTokenExpired  = fmt.Errorf("token expired")                      // 403
	ErrNotFound      = fmt.Errorf("user not found")                     // 403
	ErrGenerateToken = fmt.Errorf("generate token failed")              // 403
)
