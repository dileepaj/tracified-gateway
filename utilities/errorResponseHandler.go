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

//400 Bad Request: The request sent to the server was malformed or invalid.
func BadRequestResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "Bad request",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//404 Not Found: The requested resource could not be found on the server.
func NotFound(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "Not found",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//500 Internal Server Error: The server encountered an error while processing the request.
func InternalError(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "Internal server error",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//500 Internal Server Error: The server encountered a Database connection issue.
func DBConnectionIssue(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "Database connection issue",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//403 Forbidden: The server refuses to fulfill the request, even if authentication credentials were provided.
func ForbiddenRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusForbidden)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusForbidden,
		Error:   "Request forbidden",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//401 Unauthorized: The user does not have permission to access the requested resource.
func UnauthorizedRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnauthorized)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "Request Unauthorized",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//503 Service Unavailable: The server is currently unavailable or overloaded and cannot handle the request.
func ServiceUnavailable(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusServiceUnavailable)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusServiceUnavailable,
		Error:   "Requested service unavailable",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//429 Too Many Requests : user has sent too many requests in a given amount of time
func TooManyRequests(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusTooManyRequests)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusTooManyRequests,
		Error:   "Too many requests",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//204 No Content : Server has successfully processed the request, but there is no response body to return to the client
func NoContent(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNoContent)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusNoContent,
		Error:   "No Content",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//405 Method Not Allowed: The HTTP method used in the request is not allowed
func MethodNotAllowed(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusMethodNotAllowed,
		Error:   "Method not allowed",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//406 Not Acceptable: The server cannot respond with the content requested in the Accept headers of the request.
func NotAcceptable(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusNotAcceptable)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusNotAcceptable,
		Error:   "Not acceptable",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}

//504 Gateway Timeout: The server did not receive a timely response from another server while attempting to fulfill the request.
func GatewayTimeOut(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusGatewayTimeout)
	response := ErrorResponseBody{
		Message: message,
		Status:  http.StatusGatewayTimeout,
		Error:   "Gateway time out",
	}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		ErrorLogger.Println(err)
	}
}
