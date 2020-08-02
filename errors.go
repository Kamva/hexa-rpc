package hrpc

import (
	"github.com/kamva/hexa"
	"net/http"
)

//--------------------------------
// hrpc errors
//--------------------------------

// Error code description:
// hrpc = hrpc project (package or project name)
// u = unknown errors section (identify some part in application)
// e = Error (type of code : error|response|...)
// 0 = error number zero (id of code in that part and type)

//--------------------------------
// Unknown errors
//--------------------------------

var (
	ErrUnknownError = hexa.NewError(http.StatusInternalServerError, "hrpc.u.e.0", hexa.ErrKeyInternalError, nil)
)
