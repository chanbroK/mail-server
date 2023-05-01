package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Mail  email info
type Mail struct {
	From                *mail.Address
	To                  []*mail.Address
	Subject             string
	Body                string
	Attachment          []string
	boundaryAlternative string
	boundaryMixed       string
}

func (m *Mail) initBoundary() {
	m.boundaryAlternative = "a" + MakeBoundary()
	m.boundaryMixed = "m" + MakeBoundary()
}

func (m *Mail) getHeader() (string, error) {
	to := ""
	for i := range m.To {
		to += m.To[i].String() + ", "
	}
	header := ""
	header += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	header += fmt.Sprintf("Subject: =?utf-8?B?%s?=\r\n", base64.StdEncoding.EncodeToString(bytes.NewBufferString(m.Subject).Bytes()))
	header += fmt.Sprintf("From: %s\r\n", m.From.String())
	header += fmt.Sprintf("To: %s\r\n", to[:len(to)-2])
	header += fmt.Sprintf("MIME-Version: %s\r\n", "1.0")
	return header, nil
}

func (m *Mail) getBody() (string, error) {
	body := ""
	if len(m.Attachment) != 0 {
		body += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", m.boundaryMixed)
		body += fmt.Sprintf("--%s\r\n", m.boundaryMixed)
	}

	content, err := ChunkSplit(base64.StdEncoding.EncodeToString(bytes.NewBufferString(m.Body).Bytes()))
	if err != nil {
		return "", nil
	}

	// email body content type 및 인코딩 설정
	body += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", m.boundaryAlternative)
	body += fmt.Sprintf("--%s\r\n", m.boundaryAlternative)

	// plain text 메일 호환
	body += fmt.Sprintf("Content-Type: text/plain; charset=utf-8\r\n")
	body += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n")
	body += fmt.Sprintf("%s\r\n\r\n", content)
	body += fmt.Sprintf("--%s\r\n", m.boundaryAlternative)
	// html 메일 호환
	body += fmt.Sprintf("Content-Type: text/html; charset=utf-8\r\n")
	body += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n")
	body += fmt.Sprintf("%s\r\n\r\n", content)
	body += fmt.Sprintf("--%s--\r\n", m.boundaryAlternative)

	if len(m.Attachment) != 0 {
		for _, s := range m.Attachment {
			// TODO 서버의 file system 대신 cloud storage에서 파일 읽도록 수정
			b, err := os.ReadFile(s)
			if err != nil {
				return "", err
			}
			name := filepath.Base(s)
			data, err := ChunkSplit(base64.StdEncoding.EncodeToString(b))
			if err != nil {
				return "", err
			}
			body += fmt.Sprintf("--%s\r\n", m.boundaryMixed)
			body += fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", name)
			body += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n")
			body += fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", name)
			body += fmt.Sprintf("%s\r\n\r\n", data)
		}
		body += fmt.Sprintf("--%s--\r\n", m.boundaryMixed)
	}
	return body, nil
}

// GetReader return  io.Reader of mail message.
func (m *Mail) GetReader() (io.Reader, error) {
	m.initBoundary()
	header, err := m.getHeader()
	if err != nil {
		return nil, err
	}
	body, err := m.getBody()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(header + body), nil
}
