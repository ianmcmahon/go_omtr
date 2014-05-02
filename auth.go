package omtr

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"time"
)

type OmnitureClient struct {
	username      string
	shared_secret string
}

func New(username, shared_secret string) *OmnitureClient {
	omcl := &OmnitureClient{username, shared_secret}

	return omcl
}

func md5_hex(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func sha_64(s string) string {
	h := sha1.New()
	enc := base64.StdEncoding
	io.WriteString(h, s)

	bytes := h.Sum(nil)

	return enc.EncodeToString(bytes)
}

func (omcl *OmnitureClient) get_header() string {
	enc := base64.StdEncoding

	t := time.Now().In(time.UTC)

	md5nonce := md5_hex(t.String())
	base64nonce := enc.EncodeToString([]byte(md5nonce))
	createdDate := t.Format("2006-01-02T15:04:05Z")
	password_64 := sha_64(fmt.Sprintf("%s%s%s", md5nonce, createdDate, omcl.shared_secret))

	return fmt.Sprintf("UsernameToken Username=\"%s\", PasswordDigest=\"%s\", Nonce=\"%s\", Created=\"%s\"", omcl.username, password_64, base64nonce, createdDate)
}
