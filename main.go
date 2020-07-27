package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SetDebug() bool {
	envValue := os.Getenv("GO_RPXY_DEBUG")
	debug, err := strconv.ParseBool(envValue)
	if err != nil {
		debug = false
	}
	return debug
}

func AcceptMode() bool {
	envValue := os.Getenv("GO_PROXY_MODE")
	acceptMode, err := strconv.ParseBool(envValue)
	if err != nil {
		acceptMode = false
	}
	return acceptMode
}

func GetDomains() []*regexp.Regexp {
	envValue := os.Getenv("GO_PROXY_DOMAINS")
	domains := strings.Fields(envValue)
	domainRegexps := make([]*regexp.Regexp, len(domains))
	for idx, domain := range domains {
		domainRegexps[idx] = regexp.MustCompile(domain)
	}
	return domainRegexps
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = SetDebug()
	domains := GetDomains()

	var reqCondition goproxy.ReqConditionFunc
	if AcceptMode() {
		goproxy.SrcIpIs()
		reqCondition = goproxy.ReqHostMatches(domains...)
	} else {
		reqCondition = goproxy.Not(goproxy.ReqHostMatches(domains...))
	}

	proxy.OnRequest(reqCondition).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx)(*http.Request, *http.Response) {
			return r,goproxy.NewResponse(r,
				goproxy.ContentTypeText,http.StatusForbidden,
				"Don't waste your time!")
		})
	proxy.OnRequest(goproxy.ReqHostMatches(domains...)).HandleConnect(goproxy.AlwaysReject)
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
