package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ResponseBody struct {
	Protocols   []string `json:"protocols"`
	OS          string   `json:"os"`
	Arch        string   `json:"arch"`
	Filename    string   `json:"filename"`
	DownloadURL string   `json:"download_url"`
	SHA256Sum   string   `json:"shasum"`

	SHA256SumsURL          string `json:"shasums_url"`
	SHA256SumsSignatureURL string `json:"shasums_signature_url"`

	SigningKeys SigningKeyList `json:"signing_keys"`
}
type SigningKeyList struct {
	GPGPublicKeys []*SigningKey `json:"gpg_public_keys"`
}
type SigningKey struct {
	ASCIIArmor     string `json:"ascii_armor"`
	TrustSignature string `json:"trust_signature"`
}

func startServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := httpClient.Get(originProviderHost + r.URL.String())
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer resp.Body.Close()

		if strings.Contains(r.URL.String(), "/v1/providers/") && strings.Contains(r.URL.String(), "download") {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			data := ResponseBody{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			provider := strings.Join(strings.Split(r.URL.String(), "/")[2:4], "/")
			data.DownloadURL = fmt.Sprintf("%s/%s/%s", privateProviderHost, provider, data.Filename)
			suffix := fmt.Sprintf("%s_%s.zip", data.OS, data.Arch)
			data.SHA256SumsURL = fmt.Sprintf("%s/%s/%s", privateProviderHost, provider, strings.ReplaceAll(data.Filename, suffix, "SHA256SUMS"))
			data.SHA256SumsSignatureURL = fmt.Sprintf("%s/%s/%s", privateProviderHost, provider, strings.ReplaceAll(data.Filename, suffix, "SHA256SUMS.sig"))

			js, _ := json.Marshal(data)
			w.Write(js)

		} else {
			io.Copy(w, resp.Body)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("server start")
	http.ListenAndServe(":80", nil)
}
