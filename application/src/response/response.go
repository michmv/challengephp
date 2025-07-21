package response

import (
	"challengephp/lib"
	"net/http"
)

type Response interface {
	Render(w http.ResponseWriter) lib.Error
}
