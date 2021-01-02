package handler

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type IPAcessMiddleware struct {
	IPs  []net.IP
	IPNs []net.IPNet
}

func (ipm IPAcessMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if doIPAcess(req, ipm.IPs, ipm.IPNs) {
				fmt.Println("ip acces successful")
				nextHandler.ServeHTTP(resp, req)
			} else {
				fmt.Println("ip acces faulure")
			}
		})
	}
}

type ClientCertAccessMiddleware struct {
	IssuerCN []string
}

func (cca ClientCertAccessMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			if doCertAccess(req, cca.IssuerCN) {
				fmt.Println("cert acces successful")
				nextHandler.ServeHTTP(resp, req)
			} else {
				fmt.Println("cert acces faulure")
			}
		})
	}
}

func doIPAcess(req *http.Request, ips []net.IP, ipns []net.IPNet) bool {
	fmt.Println("ips", ips)
	fmt.Println("ipns", ipns)

	cIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Println("Error getting Ip:", err.Error())
		return false
	}

	fmt.Println("cIP", cIP)

	clientIP := net.ParseIP(cIP)
	fmt.Println("clientIP", clientIP)
	if clientIP == nil {
		fmt.Println("could not get ip")
	}

	for _, ip := range ips {
		if ip.Equal(clientIP) {
			fmt.Println("Matched IP", clientIP)
			return true
		}
	}

	for _, ipn := range ipns {
		if ipn.Contains(clientIP) {
			fmt.Println("Contains IP", clientIP)
			return true
		}
	}

	return false
}

func doCertAccess(req *http.Request, issuerCN []string) bool {
	certs := req.TLS.PeerCertificates
	if len(certs) == 0 {
		return false
	}
	cn := certs[0].Issuer.CommonName
	fmt.Printf("Found Common Name: %s\n", cn)
	for _, icn := range issuerCN {
		if strings.EqualFold(icn, cn) {
			fmt.Printf("%s matches %s\n", icn, cn)
			return true
		}
		fmt.Printf("%s doesn not matche %s\n", icn, cn)
	}
	fmt.Println("No match found for ", cn)
	return false
}
