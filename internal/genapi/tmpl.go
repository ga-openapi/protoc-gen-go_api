package genapi

import (
	"bytes"
	"log"
	"text/template"
)

var frame = `// Code generated by protoc-gen-go_api(github.com/dev-openapi/protoc-gen-go_api version={{ .Version }}). DO NOT EDIT.
// source: {{ .Source }}

package {{ .GoPackage }}

import (
	context "context"
	fmt "fmt"
	io "io"
	json "encoding/json"
	bytes "bytes"
	http "net/http"
	strings "strings"
	url "net/url"
	multipart "mime/multipart"
)
// Reference imports to suppress errors if they are not otherwise used.
var _ = context.Background
var _ = http.NewRequest
var _ = io.Copy
var _ = bytes.Compare
var _ = json.Marshal
var _ = strings.Compare
var _ = fmt.Errorf
var _ = url.Parse
var _ = multipart.ErrMessageTooLarge

{{ range .Services }}
// Client API for {{ .ServName }} service

type {{ .ServName }}Service interface {
{{- range .Methods }}
	// {{ .MethName }} {{ .Comment }}
	{{ .MethName }}(ctx context.Context, in *{{ .ReqTyp }}, opts ...Option) (*{{ .ResTyp }}, error)
{{- end }}
}

type {{ unexport .ServName }}Service struct {
	// opts
	opts *Options
}

func New{{ .ServName }}Service(opts ...Option) {{ .ServName }}Service {
	opt := newOptions(opts...)
	if len(opt.addr) <= 0 {
		opt.addr = "https://{{ .PkgName }}"
	}
	return &{{ unexport .ServName }}Service {
		opts: opt,
	}
}

{{ range .Methods }}
func (c *{{ unexport .ServName }}Service) {{ .MethName }}(ctx context.Context, in *{{ .ReqTyp }}, opts ...Option) (*{{ .ResTyp }}, error) {
	var res {{ .ResTyp }}
	{{ .ReqCode | html }}
}
{{ end -}}

{{ end -}}
`

var requestCode = `// options
	opt := buildOptions(c.opts, opts...)
	headers := make(map[string]string)
	// route
	{{ .RouteCode }}
	// body
	{{ .BodyCode }}
	{{- if eq .BodyCode "" -}}
	req, err := http.NewRequest("{{ .Verb }}", rawURL, nil)
	{{- else }}
	req, err := http.NewRequest("{{ .Verb }}", rawURL, body)
	{{- end }}
	if err != nil {
		return nil, err
	}
	{{ if ne .QueryCode "" }}
	params := req.URL.Query()
	{{ .QueryCode | html }}	
	req.URL.RawQuery = params.Encode()
	{{ end }}
	// header
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := opt.DoRequest(ctx, opt.client, req)
	if err != nil {
		return nil, err
	}
	err = opt.DoResponse(ctx, resp, &res)
	return &res, err 
`

var bodyFormCode = `bodyForms := url.Values{} 
	{{ .Body }}
	body := strings.NewReader(bodyForms.Encode())
	headers["Content-Type"] = "application/x-www-form-urlencoded"
`

var bodyMultiCode = `body := new(bytes.Buffer)
	bodyForms := multipart.NewWriter(body) 
	{{ .Body }}
	defer func() { _ =  bodyForms.Close() } ()
	headers["Content-Type"] = "multipart/form-data"
`

var bodyJsonCode = `bs, err := json.Marshal({{ .Body | html }})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	headers["Content-Type"] = "application/json"
`

var bodyByteCode = `body := bytes.NewReader({{ .Body }})
	headers["Content-Type"] = "application/json"
`

func buildFrame(data *FileData) (string, error) {
	frm, err := template.New("frame_tmpl").Funcs(fn).Parse(frame)
	if err != nil {
		log.Println("parse frame template err: ", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = frm.Execute(bs, data)
	if err != nil {
		log.Println("execuete frame template err: ", err)
		return "", err
	}
	return bs.String(), nil
}

func buildRequestCode(data *CodeData) (string, error) {
	rct, err := template.New("request_code_tmpl").Funcs(fn).Parse(requestCode)
	if err != nil {
		log.Println("parse request code template err:", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = rct.Execute(bs, data)
	if err != nil {
		log.Println("execute request code template err: ", err)
		return "", err
	}
	return bs.String(), nil
}

func buildBodyFormCode(body string) (string, error) {
	bft, err := template.New("body_form_tmpl").Funcs(fn).Parse(bodyFormCode)
	if err != nil {
		log.Println("parse form code template err: ", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = bft.Execute(bs, map[string]string{
		"Body": body,
	})
	if err != nil {
		log.Println("execute form code template err: ", err)
		return "", err
	}
	return bs.String(), nil
}

func buildBodyMultiCode(body string) (string, error) {
	bmt, err := template.New("body_multi_tmpl").Funcs(fn).Parse(bodyMultiCode)
	if err != nil {
		log.Println("parse multi code template err: ", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = bmt.Execute(bs, map[string]string{
		"Body": body,
	})
	if err != nil {
		log.Println("execute multi code template err: ", err)
		return "", err
	}
	return bs.String(), nil
}

func buildBodyJsonCode(body string) (string, error) {
	bjt, err := template.New("body_json_tmpl").Funcs(fn).Parse(bodyJsonCode)
	if err != nil {
		log.Println("parse json code template err: ", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = bjt.Execute(bs, map[string]string{
		"Body": body,
	})
	if err != nil {
		log.Println("execute json code template err: ", err)
		return "", err
	}
	return bs.String(), nil

}

func buildBodyByteCode(body string) (string, error) {
	bbt, err := template.New("body_byte_tmpl").Funcs(fn).Parse(bodyByteCode)
	if err != nil {
		log.Println("parse byte code template err: ", err)
		return "", err
	}
	bs := new(bytes.Buffer)
	err = bbt.Execute(bs, map[string]string{
		"Body": body,
	})
	if err != nil {
		log.Println("execute byte code template err: ", err)
		return "", err
	}
	return bs.String(), nil
}
