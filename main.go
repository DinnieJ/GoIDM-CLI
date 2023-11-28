package main

import (
	"bufio"
	"fmt"
	"installer/downloader"
	ihttp "installer/http"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func main() {
	other()
	return
	// testPB()
	var url string = "http://ipv4.download.thinkbroadband.com/5MB.zip"
	var metadata, _ = ihttp.GetMetadata(url)
	fmt.Printf("Filename:%s\nDownload Size: %.2f MB(s)\n", metadata.FileName, float64(metadata.ContentLength)/(1024*1024))

	if metadata.SupportPartial {
		fmt.Println("Partial supported, start download using multiple thread")
	} else {
		fmt.Println("Server not support partial download, start download single thread")
		downloader.SingleThreadDownload(&metadata)
	}
}

func other() {
	var listener, err = net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()
	mustCopy(listener, os.Stdout)
	// for {
	// 	var conn, err = listener.Accept()
	// 	if err != nil {
	// 		log.Print(err)
	// 		continue
	// 	}

	// 	go handleConnection(conn)

	// }
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(1)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var input = bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Fprintln(conn, input.Text(), 1*time.Second)
	}
}

func testPB() {
	var wg sync.WaitGroup
	// passed wg will be accounted at p.Wait() call
	p := mpb.New(mpb.WithWaitGroup(&wg))
	total, numBars := 100, 3
	wg.Add(numBars)

	for i := 0; i < numBars; i++ {
		name := fmt.Sprintf("Bar#%d ", i)
		bar := p.AddBar(int64(total),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(name),
				// decor.DSyncWidth bit enables column width synchronization
				decor.AverageETA(decor.ET_STYLE_HHMMSS),
			),
			mpb.BarStyle("⟪▊▌ ⟫"),
			mpb.AppendDecorators(
				// replace ETA decorator with "done" message, OnComplete event
				decor.OnComplete(
					// ETA decorator with ewma age of 30
					decor.Percentage(decor.WCSyncWidth), "Done",
				),
			),
		)
		// simulating some work
		go func() {
			defer wg.Done()
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			max := 100 * time.Millisecond
			for i := 0; i < total; i++ {
				// start variable is solely for EWMA calculation
				// EWMA's unit of measure is an iteration's duration
				// start := time.Now()
				time.Sleep(time.Duration(rng.Intn(10)+1) * max / 10)
				// we need to call EwmaIncrement to fulfill ewma decorator's contract
				bar.IncrBy(1)
			}
		}()
	}
	// wait for passed wg and for all bars to complete and flush
	p.Wait()
	log.Fatal("Done")
}
