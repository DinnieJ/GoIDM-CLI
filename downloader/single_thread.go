package downloader

import (
	"fmt"
	ihttp "installer/http"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

func SingleThreadDownload(metadata *ihttp.Metadata) {

	var req, _ = http.NewRequest("GET", metadata.Url, nil)
	req.Header.Add("User-Agent", ihttp.USER_AGENT)
	var resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()
	var tmpFilePath = GetTmpFilePath(metadata.FileName)
	var f, _ = os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		int64(metadata.ContentLength),
		tmpFilePath,
	)
	var buf = make([]byte, metadata.ContentLength)
	fmt.Println(resp.Body)
	var total = 0
	for {
		var n, err = resp.Body.Read(buf)
		total += n
		fmt.Println(total)
		if err == io.EOF {
			break
		}
	}
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	os.Rename(tmpFilePath, fmt.Sprintf("./%s", metadata.FileName))
}
