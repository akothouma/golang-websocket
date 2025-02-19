package main

import (
	"fmt"
	"net/http"
)

func (dep *Dependencies)HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to Forum")
}
