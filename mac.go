package main

import (
	"net/http"
	"sniffer/monitors/proxy3"
)

func main() {
	// proxy := goproxy.NewProxyHttpServer()
	// proxy.Verbose = true
	//
	//proxy.OnRequest(goproxy.DstHostIs("www.reddit.com")).DoFunc(
	//	func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//		print("hello_world")
	//		if h, _, _ := time.Now().Clock(); h >= 8 && h <= 17 {
	//			return r, goproxy.NewResponse(r,
	//				goproxy.ContentTypeText, http.StatusForbidden,
	//				"Don't waste your time!")
	//		}
	//		return r, nil
	//	})

	myHandler := proxy3.NewProxy()
	println("i am started")
	http.ListenAndServe(":8080", myHandler)

	//log.Fatal(http.ListenAndServe(":8080", proxy))

}

//import (
//	"fmt"
//	//"sniffer/monitors/proxy"
//	"sniffer/monitors/proxy3"
//)
//
//func main() {
//	// go network.NetworkMonitor2()
//	// printer.MacPrinter()
//	// keyloger.MacKeyloger()
//	// keyloger.Keyloger()
//	fmt.Println("service start")
//	proxy3.Proxy3()
//	//fs.FSWatch()
//}
