package utilities

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponseBody struct {
	Message string
	Status  int
	Error   string
}

var ErrorLogger *log.Logger

func ErrorResponse(w http.ResponseWriter, message string, statusCode int, errorCode string) {
	w.WriteHeader(statusCode)
	response := ErrorResponseBody{
		Message: message,
		Status:  statusCode,
		Error:   errorCode,
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//Handle error responses (Parameters - http.ResponseWriter, message, error code)
//
//Supports
//204, 400, 401, 403, 404, 405, 406, 429, 500, 503, 504.
func HandleError(w http.ResponseWriter, message string, errorCode int) {
	switch errorCode {
	case http.StatusBadRequest:
		ErrorResponse(w, message, http.StatusBadRequest, "Bad request")
	case http.StatusNotFound:
		ErrorResponse(w, message, http.StatusNotFound, "Not found")
	case http.StatusInternalServerError:
		ErrorResponse(w, message, http.StatusInternalServerError, "Internal server error")
	case http.StatusForbidden:
		ErrorResponse(w, message, http.StatusForbidden, "Request forbidden")
	case http.StatusUnauthorized:
		ErrorResponse(w, message, http.StatusUnauthorized, "Request Unauthorized")
	case http.StatusServiceUnavailable:
		ErrorResponse(w, message, http.StatusServiceUnavailable, "Requested service unavailable")
	case http.StatusTooManyRequests:
		ErrorResponse(w, message, http.StatusTooManyRequests, "Too many requests")
	case http.StatusNoContent:
		ErrorResponse(w, message, http.StatusNoContent, "No Content")
	case http.StatusMethodNotAllowed:
		ErrorResponse(w, message, http.StatusMethodNotAllowed, "Method not allowed")
	case http.StatusNotAcceptable:
		ErrorResponse(w, message, http.StatusNotAcceptable, "Not acceptable")
	case http.StatusGatewayTimeout:
		ErrorResponse(w, message, http.StatusGatewayTimeout, "Gateway time out")
	default:
		ErrorResponse(w, message, http.StatusInternalServerError, "Unknown error")
	}
}
