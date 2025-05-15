package controllers

import (
	"fmt"
	"net/http"
	"os"
)

func HeartBeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong from "+os.Getenv("APP_NAME")+"["+os.Getenv("COMMIT")+"]\n")
}
