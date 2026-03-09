package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kasteion/httpfromtcp/internal/headers"
	"github.com/kasteion/httpfromtcp/internal/request"
	"github.com/kasteion/httpfromtcp/internal/response"
	"github.com/kasteion/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handler400(w, req)
	case "/myproblem":
		handler500(w, req)
	case "/video":
		handlerVideo(w, req)
	default:
		handler200(w, req)
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	message := `
	<html>
	<head>
		<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
	</body>
	</html>`
	w.WriteStatusLine(response.StatusCodeBadRequest)
	h := response.GetDefaultHeaders(len(message))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(message))
}

func handler500(w *response.Writer, _ *request.Request) {
	message := `
	<html>
	<head>
		<title>500 Internal Server Error</title>
	</head>
	<body>
		<h1>Internal Server Error</h1>
		<p>Okay, you know what? This one is on me.</p>
	</body>
	</html>`
	w.WriteStatusLine(response.StatusCodeInternalServerError)
	h := response.GetDefaultHeaders(len(message))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(message))
}

func handler200(w *response.Writer, _ *request.Request) {
	message := `
	<html>
	<head>
		<title>200 OK</title>
	</head>
	<body>
		<h1>Success!</h1>
		<p>Your request was an absolute banger.</p>
	</body>
	</html>`
	w.WriteStatusLine(response.StatusCodeOK)
	h := response.GetDefaultHeaders(len(message))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(message))
}

func proxyHandler(w *response.Writer, r *request.Request) {
	// target := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
	
	// resp, err := http.Get("https://httpbin.org" + target)
	resp, err := http.Get("https://httpbin.org/html")
	if err != nil {
		handler500(w, r)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeOK)
	h := response.GetDefaultHeaders(0)
	h.Remove("Content-Length")
	h.Set("Transfer-Encoding", "chunked")
	h.Set("X-Content-SHA256", "X-Content-Sha256, X-Content-Length")
	w.WriteHeaders(h)

	buf := make([]byte, 1024)
	fullBody := []byte{}
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			fullBody = append(fullBody, buf[:n]...)
			w.WriteChunkedBody(buf[:n])
		}
		if err != nil {
			if err == io.EOF{
				w.WriteChunkedBodyDone()
				fullBodyHash := sha256.Sum256(fullBody)
				t := headers.NewHeaders()
				t.Set("X-Content-Sha256", fmt.Sprintf("%x", fullBodyHash)) 
				t.Set("X-Content-Length", strconv.Itoa(len(fullBody)))
				w.WriteTrailers(t)
				break
			}
			return
		}
	}
}

func handlerVideo(w *response.Writer, r *request.Request) {
	data, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		handler500(w, r)
		return
	}

	w.WriteStatusLine(response.StatusCodeOK)
	h := response.GetDefaultHeaders(len(data))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(data)
}
