package main

import (
	"flag"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"9fans.net/go/plumb"
        "9fans.net/go/plan9"
)

var (
	plumbfile = flag.String("p", "send", "write the message to plumbfile (default send)")
	attributes = flag.String("a", "", "set the attr field of the message (default is empty), expects key=value")
	source = flag.String("s", "", "set the src field of the message (default is store)")
	destination = flag.String("d", "", "set the dst filed of the message (default is store)")
        typefield = flag.String("t", "", "override the type message sent to the plumber")
	directory = flag.String("w", "", "set the wdir field of the message (default is current directory)")
        input = flag.Bool("i", false, "take the data from standard input rather than the argument strings")
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
	var fd io.Writer
	var err error
	
	switch runtime.GOOS {
	case "plan9":
		fd, err = os.OpenFile("/mnt/plumb/send", os.O_WRONLY, 0644)
	default:
		fd, err = plumb.Open(*plumbfile, plan9.OWRITE)
	}
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

func newStoreMsg(mediaType, wdir, data string, attr *plumb.Attribute) *storeMsg {
	sf := &storeMsg{
		src: os.Args[0],
		dst: "",
		wdir: wdir,
		msgtype: mediaType,
		attr: attr,
		ndata: len(data),
		data: data,
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

func contentTypeUrl(data string) (string, error) {
	// We read in 512 bytes 
	buf := make([]byte, 512)
	u, err := url.ParseRequestURI(data)
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

func contentTypeFile(data string) (string, error) {
	buf := make([]byte, 512)
	fd, err := os.Open(data)
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

func content(data string) (string, error) {
        if *typefield != "" {
            return *typefield, nil
        }
 
	if _, err := os.Stat(data); os.IsNotExist(err)  {
		return contentTypeUrl(data)
	}
	return contentTypeFile(data)
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

func readInput() string {
        data := new(strings.Builder)
        if *input {
		// readInput on stdin will come with a fancy newline
		// return a slice without the last character 
		io.Copy(data, os.Stdin)
		return data.String()[:data.Len()-1]
        } else {
		data.WriteString(flag.Arg(0))
        }

	return data.String()
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

        data := readInput()
	ct, err := content(data)
	if err != nil {
		log.Fatal(err)
	}
	mediaType, params := getMediaType(ct)
	storeMsg := newStoreMsg(mediaType, wdir, data, params)
	err = storeMsg.send()
	if err != nil {
		log.Fatal(err)
	}
}
