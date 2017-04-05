package api

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/base64"
	"fmt"
)

func Rc4Decode(rawkey, src []byte) (dst []byte, err error) {
	var src_base64 []byte
	src_base64, err = base64.StdEncoding.DecodeString(string(src))
	if err != nil {
		return
	}

	mmd5 := md5.New()
	mmd5.Write(rawkey)
	key := fmt.Sprintf("%x", mmd5.Sum(nil))

	var c *rc4.Cipher
	c, err = rc4.NewCipher([]byte(key))
	if err != nil {
		return
	}

	dst = make([]byte, len(src_base64))
	c.XORKeyStream(dst, src_base64)

	return
}

func Rc4DecodeString(rawkey, src string) (string, error) {
	src_base64, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return "", nil
	}

	mmd5 := md5.New()
	mmd5.Write([]byte(rawkey))
	key := fmt.Sprintf("%x", mmd5.Sum(nil))

	var c *rc4.Cipher
	c, err = rc4.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	var dst []byte = make([]byte, len(src_base64))
	c.XORKeyStream(dst, src_base64)

	return string(dst), nil
}
