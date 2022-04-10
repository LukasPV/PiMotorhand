package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"os"
	_ "os/exec"
	"strings"
	"sync"
)

type myHandler struct {
	mu sync.Mutex // guards pin1, pin2

	pin11 rpio.Pin
	pin12 rpio.Pin
	pin21 rpio.Pin
	pin22 rpio.Pin
}

func NewHandler() *myHandler {
	h := &myHandler{
		pin11: rpio.Pin(23),
		pin12: rpio.Pin(24),
		pin21: rpio.Pin(5),
		pin22: rpio.Pin(6),
	}

	h.pin11.Output()
	h.pin11.Low()
	h.pin12.Output()
	h.pin12.Low()
	h.pin21.Output()
	h.pin21.Low()
	h.pin22.Output()
	h.pin22.Low()

	return h
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.anaus(r.URL.Path)

	if r.URL.Path == "/version" {
		fmt.Fprintf(w, "123 \n")
	}

}
func main() {
	fmt.Println("opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	h := NewHandler()
	http.Handle("/version", h)

	http.Handle("/Forward/Motor1", h)
	http.Handle("/Forward/Motor2", h)
	http.Handle("/Forward/Both", h)

	http.Handle("/Back/Motor1", h)
	http.Handle("/Back/Motor2", h)
	http.Handle("/Back/Both", h)

	http.Handle("/Stop/Motor1", h)
	http.Handle("/Stop/Motor2", h)
	http.Handle("/Stop/Both", h)

	path := "/static"
	directory := os.Getenv("ASSET_ROOT")
	if len(directory) == 0 {
		panic("EnvNotSet")
	}
	http.Handle("/",
		http.StripPrefix(strings.TrimRight(path, "/"), http.FileServer(http.Dir(directory))))

	bindaddr := ":8001"
	fmt.Printf("serving on %s\n", bindaddr)
	err = http.ListenAndServe(bindaddr, nil)
	if err != nil {
		panic(err)
	}
}
func (h *myHandler) anaus(s string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Println(s)
	if s == "/Forward/Motor1" {
		h.pin12.Low()
		h.pin11.High()
	}
	if s == "/Forward/Motor2" {
		h.pin22.Low()
		h.pin21.High()
	}
	if s == "/Forward/Both" {
		h.pin12.Low()
		h.pin11.High()
		h.pin22.Low()
		h.pin21.High()
	}

	if s == "/Back/Motor1" {
		h.pin11.Low()
		h.pin12.High()
	}
	if s == "/Back/Motor2" {
		h.pin21.Low()
		h.pin22.High()
	}
	if s == "/Back/Both" {
		h.pin11.Low()
		h.pin12.High()
		h.pin21.Low()
		h.pin22.High()
	}

	if s == "/Stop/Motor1" {
		h.pin11.Low()
		h.pin12.Low()
	}
	if s == "/Stop/Motor2" {
		h.pin21.Low()
		h.pin22.Low()
	}
	if s == "/Stop/Both" {
		h.pin11.Low()
		h.pin12.Low()
		h.pin21.Low()
		h.pin22.Low()
	}

}
