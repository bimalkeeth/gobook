package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8088")
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan struct{})

	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}

	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
