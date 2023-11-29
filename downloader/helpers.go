package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const TMP_FILE_PATH string = "/tmp/%s.file"

func GetTmpFilePath(fileName string, salt string) string {
	var hasher = sha1.New()
	hasher.Write([]byte(fileName + salt))
	var hashFilename = hex.EncodeToString(hasher.Sum(nil))

	return fmt.Sprintf(TMP_FILE_PATH, hashFilename)
}
