// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"math/rand"
// 	"net"
// 	"net/http"
// 	"time"
// )

// func main() {
// 	rand.Seed(time.Now().UnixNano())

// 	fwd := &forwarder{"api", 8080}
// 	http.Handle("/words/", http.StripPrefix("/words", fwd))
// 	http.Handle("/", http.FileServer(http.Dir("static")))

// 	fmt.Println("Listening on port 80")
// 	http.ListenAndServe(":80", nil)
// }

// type forwarder struct {
// 	host string
// 	port int
// }

// func (f *forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	addrs, err := net.LookupHost(f.host)
// 	if err != nil {
// 		log.Println("Error", err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	log.Printf("%s %d available ips: %v", r.URL.Path, len(addrs), addrs)
// 	ip := addrs[rand.Intn(len(addrs))]
// 	log.Printf("%s I choose %s", r.URL.Path, ip)

// 	url := fmt.Sprintf("http://%s:%d%s", ip, f.port, r.URL.Path)
// 	log.Printf("%s Calling %s", r.URL.Path, url)

// 	if err = copy(url, ip, w); err != nil {
// 		log.Println("Error", err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// }

// func copy(url, ip string, w http.ResponseWriter) error {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}

// 	for header, values := range resp.Header {
// 		for _, value := range values {
// 			w.Header().Add(header, value)
// 		}
// 	}
// 	w.Header().Set("source", ip)

// 	buf, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = w.Write(buf)
// 	return err
// }

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// ServiceForwarder knows where the backend service is.
type ServiceForwarder struct {
	host string
	port int
}

func main() {
	// Initialize the random number generator.
	rand.Seed(time.Now().UnixNano())

	// Requests to /words/* are forwarded to the API service.
	apiForwarder := &ServiceForwarder{
		host: "api",
		port: 8080,
	}

	http.Handle("/words/", http.StripPrefix("/words", apiForwarder))

	// Serve files from the "static" directory.
	http.Handle("/", http.FileServer(http.Dir("static")))

	fmt.Println("Listening on port 8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}

// ServeHTTP is called automatically whenever a request arrives.
func (f *ServiceForwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Find all IP addresses for the backend service.
	addresses, err := net.LookupHost(f.host)
	if err != nil {
		log.Println("Lookup failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s - Available IPs: %v", r.URL.Path, addresses)

	// Pick one IP address at random.
	selectedIP := addresses[rand.Intn(len(addresses))]
	log.Printf("%s - Selected: %s", r.URL.Path, selectedIP)

	// Build the backend URL.
	backendURL := fmt.Sprintf(
		"http://%s:%d%s",
		selectedIP,
		f.port,
		r.URL.Path,
	)

	log.Printf("%s - Forwarding to %s", r.URL.Path, backendURL)

	// Fetch the response from the backend and send it back to the client.
	err = forwardResponse(backendURL, selectedIP, w)
	if err != nil {
		log.Println("Forwarding failed:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// forwardResponse fetches data from the backend service and
// copies the response back to the client.
func forwardResponse(url string, sourceIP string, w http.ResponseWriter) error {

	// Send the request.
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Copy all response headers.
	for headerName, values := range response.Header {
		for _, value := range values {
			w.Header().Add(headerName, value)
		}
	}

	// Add the IP address of the backend that handled the request.
	w.Header().Set("source", sourceIP)

	// Read the response body.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Send it to the client.
	_, err = w.Write(body)
	return err
}
