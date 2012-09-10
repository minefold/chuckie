package main

import (
	"github.com/kr/s3"
	"net/http"
	"os"
	"strings"
	"time"
	"log"
	"io"
)

func main() {
	// TODO Calculate from the web request
	url  := "http://minefold-production-worlds.s3.amazonaws.com/4fdf55658e009b00010000c2/4fdf55658e009b00010000c2.1340037833.tar.gz"
	
	keys := s3.Keys{
		os.Getenv("AWS_ACCESS_KEY"),
		os.Getenv("AWS_SECRET_KEY"),
	}
	
	data := strings.NewReader("hello, world")
	
	r, _ := http.NewRequest("GET", url, data)
	r.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	r.Header.Set("X-Amz-Acl", "public-read")
	s3.Sign(r, keys)
	
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	fo, err := os.Create("output.zip")
	
	buf := make([]byte, 1024)
	
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		
		if n == 0 {
			break
		}
		
		if n2, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		} else if n2 != n {
			panic("error in writing")
		}
	}
}
