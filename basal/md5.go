package basal

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func MD5String(s string) string {
	return MD5Bytes([]byte(s))
}

func MD5Bytes(s []byte) string {
	md5Ctx := md5.New()
	md5Ctx.Write(s)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func MD5File(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return MD5Bytes(data), nil
}
