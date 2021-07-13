package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"9fans.net/go/plumb"
        "9fans.net/go/plan9"
)

var (
	plumbfile = flag.String("p", "send", "write the message to plumbfile (default send)")
	attributes = flag.String("a", "", "set the attr field of the message (default is empty), expects key=value")
	source = flag.String("s", "", "set the src field of the message (default is store)")
	destination = flag.String("d", "", "set the dst filed of the message (default is store)")
	directory = flag.String("w", "", "set the wdir field of the message (default is current directory)")
)

type storeMsg struct {
	src string
	dst string
	wdir string
	msgtype string
	attr *plumb.Attribute
	ndata int
	data string
}

func (s storeMsg) send() error {
	fd, err := plumb.Open("send", plan9.OWRITE)
	if err != nil {
		return err
	}
	message := &plumb.Message{
		Src: s.src,
		Dst: s.dst,
		Dir: s.wdir,
		Type: s.msgtype,
		Attr: s.attr,
		Data: []byte(s.data),
	}
	return message.Send(fd)
}

func newStoreMsg(mediaType, wdir, arg string, attr *plumb.Attribute) *storeMsg {
	sf := &storeMsg{
		src: os.Args[0],
		dst: "",
		wdir: wdir,
		msgtype: mediaType,
		attr: attr,
		ndata: len(arg),
		data: arg,
	}
	if *source != "" {
		sf.src = *source
	}
	if *destination != "" {
		sf.dst = *destination
	}
	return sf
}

func paramsToAttr(params map[string]string) *plumb.Attribute {
	// Attribute is a linked list - we only get one from content-type, the encoding
	attr := &plumb.Attribute{Name: "", Value: ""}
	for key, value := range params {
		if (attr.Name == "" || attr.Value == "") {
			continue
		}
		attr.Name = key
		attr.Value = value
	}
	if *attributes != "" {
		attr.Name = strings.TrimLeft(*attributes, "=")
		attr.Value = strings.TrimRight(*attributes, "=")
	}
	return attr
}

func contentTypeUrl(arg string) (string, error) {
	// We read in 512 bytes 
	buf := make([]byte, 512)
	u, err := url.ParseRequestURI(arg)
	if err != nil {
		return "text", nil
	}
	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	n, err := r.Body.Read(buf)
	if err != nil {
		return "", err
	}
	mediaType := http.DetectContentType(buf[:n])
	if (mediaType == "application/octet-stream" || mediaType == "" ) {
		mediaType := r.Header.Get("Content-type")
		if mediaType == "" {
			return "text", nil
		}
	}
	return mediaType, nil
}

func contentTypeFile(arg string) (string, error) {
	buf := make([]byte, 512)
	fd, err := os.Open(arg)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	n, err := fd.Read(buf)
	if err != nil {
		return "", err
	}
	mediaType := http.DetectContentType(buf[:n])
	if mediaType == "application/octet-stream" {
		return "text", nil
	}
	return mediaType, nil
}

func content(arg string) (string, error) {
	if _, err := os.Stat(arg); os.IsNotExist(err)  {
		return contentTypeUrl(arg)
	}
	return contentTypeFile(arg)
}

func getMediaType(ct string) (string, *plumb.Attribute) {
		if ct == "text" {
			return ct, nil
		}
		mediaType, params, err := mime.ParseMediaType(ct)
		if err != nil {
			log.Fatal(err)
		}
		return mediaType, paramsToAttr(params)
}

func main() {
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	wdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, arg := range os.Args[1:] {
		ct, err := content(arg)
		if err != nil {
			log.Fatal(err)
		}
		mediaType, params := getMediaType(ct)
		storeMsg := newStoreMsg(mediaType, wdir, arg, params)
		err = storeMsg.send()
		if err != nil {
			log.Fatal(err)
		}
	}
}
