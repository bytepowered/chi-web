package chiweb

import (
	"github.com/bytedance/sonic"
	"log/slog"
	"net/http"
)

func SendJSON(w http.ResponseWriter, statusCode int, data []byte) (int, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return w.Write(data)
}

func SendTEXT(w http.ResponseWriter, statusCode int, data []byte) (int, error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)
	return w.Write(data)
}

func SendJSONObjectOK(w http.ResponseWriter, data any) {
	SendJSONObject(w, http.StatusOK, data)
}

func SendJSONObject(w http.ResponseWriter, statusCode int, data any) {
	bytes, err := sonic.ConfigFastest.Marshal(map[string]any{
		"status": "success",
		"data":   data,
	})
	if err != nil {
		slog.Error("app: failed to marshal data: %v", err)
		SendERROR(w, http.StatusInternalServerError, "data serialized error")
	} else {
		SendJSON(w, statusCode, bytes)
	}
}

func SendERROR(w http.ResponseWriter, statusCode int, msg string) {
	bytes, _ := sonic.ConfigFastest.Marshal(map[string]any{
		"status": "error",
		"msg":    msg,
	})
	SendJSON(w, statusCode, bytes)
}
