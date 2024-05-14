package proxy3

//import (
//	"errors"
//	"flag"
//	"fmt"
//	"github.com/elazarl/goproxy"
//	"github.com/elazarl/goproxy/transport"
//	"io"
//	"log"
//	"net"
//	"net/http"
//	"net/http/httputil"
//	"os"
//	"os/signal"
//	"path"
//	"sync"
//	"time"
//)

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"path"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/transport"
)

type FileStream struct {
	path string
	f    *os.File
}

func NewFileStream(path string) *FileStream {
	return &FileStream{path, nil}
}

func (fs *FileStream) Write(b []byte) (nr int, err error) {
	if fs.f == nil {
		fs.f, err = os.Create(fs.path)
		if err != nil {
			return 0, err
		}
	}
	return fs.f.Write(b)
}

func (fs *FileStream) Close() error {
	fmt.Println("Close", fs.path)
	if fs.f == nil {
		return errors.New("FileStream was never written into")
	}
	return fs.f.Close()
}

type Meta struct {
	req      *http.Request
	resp     *http.Response
	err      error
	t        time.Time
	sess     int64
	bodyPath string
	from     string
}

func fprintf(nr *int64, err *error, w io.Writer, pat string, a ...interface{}) {
	if *err != nil {
		return
	}
	var n int
	n, *err = fmt.Fprintf(w, pat, a...)
	*nr += int64(n)
}

func write(nr *int64, err *error, w io.Writer, b []byte) {
	if *err != nil {
		return
	}
	var n int
	n, *err = w.Write(b)
	*nr += int64(n)
}

func (m *Meta) WriteTo(w io.Writer) (nr int64, err error) {
	if m.req != nil {
		fprintf(&nr, &err, w, "Type: request\r\n")
	} else if m.resp != nil {
		fprintf(&nr, &err, w, "Type: response\r\n")
	}
	fprintf(&nr, &err, w, "ReceivedAt: %v\r\n", m.t)
	fprintf(&nr, &err, w, "Session: %d\r\n", m.sess)
	fprintf(&nr, &err, w, "From: %v\r\n", m.from)
	if m.err != nil {
		// note the empty response
		fprintf(&nr, &err, w, "Error: %v\r\n\r\n\r\n\r\n", m.err)
	} else if m.req != nil {
		fprintf(&nr, &err, w, "\r\n")
		buf, err2 := httputil.DumpRequest(m.req, false)
		if err2 != nil {
			return nr, err2
		}
		write(&nr, &err, w, buf)
	} else if m.resp != nil {
		fprintf(&nr, &err, w, "\r\n")
		buf, err2 := httputil.DumpResponse(m.resp, false)
		if err2 != nil {
			return nr, err2
		}
		write(&nr, &err, w, buf)
	}
	return
}

// HttpLogger is an asynchronous HTTP request/response logger. It traces
// requests and responses headers in a "log" file in logger directory and dumps
// their bodies in files prefixed with the session identifiers.
// Close it to ensure pending items are correctly logged.
type HttpLogger struct {
	path  string
	c     chan *Meta
	errch chan error
}

func NewLogger(basepath string) (*HttpLogger, error) {
	f, err := os.Create(path.Join(basepath, "log"))
	if err != nil {
		return nil, err
	}
	logger := &HttpLogger{basepath, make(chan *Meta), make(chan error)}
	go func() {
		for m := range logger.c {
			if _, err := m.WriteTo(f); err != nil {
				log.Println("Can't write meta", err)
			}
		}
		logger.errch <- f.Close()
	}()
	return logger, nil
}

func (logger *HttpLogger) LogResp(resp *http.Response, ctx *goproxy.ProxyCtx) {
	body := path.Join(logger.path, fmt.Sprintf("%d_resp", ctx.Session))
	from := ""
	if ctx.UserData != nil {
		from = ctx.UserData.(*transport.RoundTripDetails).TCPAddr.String()
	}
	if resp == nil {
		resp = emptyResp
	} else {
		resp.Body = NewTeeReadCloser(resp.Body, NewFileStream(body))
	}
	logger.LogMeta(&Meta{
		resp: resp,
		err:  ctx.Error,
		t:    time.Now(),
		sess: ctx.Session,
		from: from})
}

var emptyResp = &http.Response{}
var emptyReq = &http.Request{}

func (logger *HttpLogger) LogReq(req *http.Request, ctx *goproxy.ProxyCtx) {
	body := path.Join(logger.path, fmt.Sprintf("%d_req", ctx.Session))
	if req == nil {
		req = emptyReq
	} else {
		req.Body = NewTeeReadCloser(req.Body, NewFileStream(body))
	}
	logger.LogMeta(&Meta{
		req:  req,
		err:  ctx.Error,
		t:    time.Now(),
		sess: ctx.Session,
		from: req.RemoteAddr})
}

func (logger *HttpLogger) LogMeta(m *Meta) {
	logger.c <- m
}

func (logger *HttpLogger) Close() error {
	close(logger.c)
	return <-logger.errch
}

// TeeReadCloser extends io.TeeReader by allowing reader and writer to be
// closed.
type TeeReadCloser struct {
	r io.Reader
	w io.WriteCloser
	c io.Closer
}

func NewTeeReadCloser(r io.ReadCloser, w io.WriteCloser) io.ReadCloser {
	return &TeeReadCloser{io.TeeReader(r, w), w, r}
}

func (t *TeeReadCloser) Read(b []byte) (int, error) {
	return t.r.Read(b)
}

// Close attempts to close the reader and write. It returns an error if both
// failed to Close.
func (t *TeeReadCloser) Close() error {
	err1 := t.c.Close()
	err2 := t.w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// stoppableListener serves stoppableConn and tracks their lifetime to notify
// when it is safe to terminate the application.
type stoppableListener struct {
	net.Listener
	sync.WaitGroup
}

type stoppableConn struct {
	net.Conn
	wg *sync.WaitGroup
}

func newStoppableListener(l net.Listener) *stoppableListener {
	return &stoppableListener{l, sync.WaitGroup{}}
}

func (sl *stoppableListener) Accept() (net.Conn, error) {
	c, err := sl.Listener.Accept()
	if err != nil {
		return c, err
	}
	sl.Add(1)
	return &stoppableConn{c, &sl.WaitGroup}, nil
}

func (sc *stoppableConn) Close() error {
	sc.wg.Done()
	return sc.Conn.Close()
}

func Proxy3() {
	verbose := flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr := flag.String("l", ":8080", "on which address should the proxy listen")
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = *verbose
	if err := os.MkdirAll("db", 0755); err != nil {
		log.Fatal("Can't create dir", err)
	}
	logger, err := NewLogger("db")
	if err != nil {
		log.Fatal("can't open log file", err)
	}
	tr := transport.Transport{Proxy: transport.ProxyFromEnvironment}
	// For every incoming request, override the RoundTripper to extract
	// connection information. Store it is session context log it after
	// handling the response.
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		ctx.RoundTripper = goproxy.RoundTripperFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (resp *http.Response, err error) {
			ctx.UserData, resp, err = tr.DetailedRoundTrip(req)
			return
		})
		logger.LogReq(req, ctx)
		return req, nil
	})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		logger.LogResp(resp, ctx)
		return resp
	})
	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("listen:", err)
	}
	sl := newStoppableListener(l)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Println("Got SIGINT exiting")
		sl.Add(1)
		sl.Close()
		logger.Close()
		sl.Done()
	}()
	log.Println("Starting Proxy")
	http.Serve(sl, proxy)
	sl.Wait()
	log.Println("All connections closed - exit")
}

//
//import (
//	"bufio"
//	"flag"
//	"github.com/elazarl/goproxy"
//	"log"
//	"net"
//	"net/http"
//	"regexp"
//)
//
//func orPanic(err error) {
//	if err != nil {
//		panic(err)
//	}
//}
//
//func Proxy3() {
//	log.Println("START PROXY!!!!!")
//
//	proxy := goproxy.NewProxyHttpServer()
//	//proxy.OnRequest().DoFunc(
//	//	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
//	//		println("hellox")
//	//		log.Println("reques", r.Host)
//	//
//	//		r.Header.Set("X-GoProxy", "yxorPoG-X")
//	//		return r, nil
//	//	})
//	// proxy.CertStore = NewCertStorage() //设置storage
//	// proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
//
//	//proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
//	//	log.Println("======")
//	//	log.Println(req.URL.String())
//	//	return req, nil
//	//})
//
//	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
//		HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
//			defer func() {
//				if e := recover(); e != nil {
//					ctx.Logf("error connecting to remote: %v", e)
//					client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
//				}
//				client.Close()
//			}()
//			log.Println("startingsuper")
//			clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
//			remote, err := net.Dial("tcp", req.URL.Host)
//			orPanic(err)
//			client.Write([]byte("HTTP/1.1 200 Ok\r\n\r\n"))
//			remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
//			for {
//				req, err := http.ReadRequest(clientBuf.Reader)
//				orPanic(err)
//				orPanic(req.Write(remoteBuf))
//				orPanic(remoteBuf.Flush())
//				resp, err := http.ReadResponse(remoteBuf.Reader, req)
//				orPanic(err)
//				orPanic(resp.Write(clientBuf.Writer))
//				orPanic(clientBuf.Flush())
//			}
//		})
//
//	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
//	// proxy.Verbose = true
//	proxy.Verbose = *verbose
//	log.Fatal(http.ListenAndServe(":8080", proxy))
//}
