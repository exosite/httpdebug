package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type HandleAll struct {
	Verbose bool
}

func (self *HandleAll) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if self.Verbose {
		log.Printf("%s %s", req.Method, req.URL.String())
	}

	reqDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("Failed to dump request: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
	}
	log.Print(string(reqDump))
	
	w.WriteHeader(http.StatusOK)
	w.Write(reqDump)
}

func main() {
	hAll := HandleAll{
		Verbose: os.Getenv("VERBOSE") == "y",
	}
	log.Fatal(http.ListenAndServe(":8080", &hAll))
}
