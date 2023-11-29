package downloader

import (
	"bytes"
	"fmt"
	c_ "installer/common"
	ihttp "installer/http"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
)

const BUF_CHUNK_SIZE = 1024 * 1024 * 5 //5MB for buffer, fuck me I have lots of RAM

func SingleThreadDownload(metadata *ihttp.Metadata) {

	var req, _ = http.NewRequest("GET", metadata.Url, nil)
	req.Header.Add("User-Agent", ihttp.USER_AGENT)
	var resp, _ = http.DefaultClient.Do(req)
	defer resp.Body.Close()
	var tmpFilePath = GetTmpFilePath(metadata.FileName, c_.RandStringRunes(8))
	var f, _ = os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()
	fmt.Println(int64(metadata.ContentLength), metadata.ContentLength)
	bar := progressbar.DefaultBytes(
		int64(metadata.ContentLength),
		tmpFilePath,
	)

	var total = 0
	for {
		var buf = make([]byte, BUF_CHUNK_SIZE)
		var n, err = resp.Body.Read(buf)
		total += n
		f.Write(buf[:n])
		io.Copy(bar, bytes.NewReader(buf[:n]))
		buf = make([]byte, BUF_CHUNK_SIZE)
		if err == io.EOF {
			break
		}
	}

	os.Rename(tmpFilePath, fmt.Sprintf("./t%s", metadata.FileName))
}
