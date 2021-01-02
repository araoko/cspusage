package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/araoko/cspusage/config"
	"github.com/araoko/cspusage/entity"
	"github.com/araoko/cspusage/handler"
	"github.com/araoko/cspusage/util"

	"github.com/gorilla/mux"
)

func getRouter(c *config.Config, logger *ServiceLogger) (*mux.Router, error) {
	db, err := entity.GetMySQLDB(c)
	if err != nil {
		return nil, err
	}
	auth := entity.GetADAuth(c)

	ips, ipns, err := getIps(c, logger)

	ipacceMW := handler.IPAcessMiddleware{
		IPs:  ips,
		IPNs: ipns,
	}

	ccAccessMw := handler.ClientCertAccessMiddleware{
		IssuerCN: c.IssuerCNs,
	}

	router := mux.NewRouter()

	router.Handle("/",
		ccAccessMw.Middleware()(ipacceMW.Middleware()(http.HandlerFunc(landing)))).Methods("GET")
	router.Handle("/export/{b64}", handler.ExportHandler{DB: db})

	router.Handle("/list", handler.CustomerListHandler{DB: db})
	router.Handle("/customermonthly", handler.CustomerMonthlyBillHandler{DB: db})
	router.Handle("/customerbillrange", handler.CustomerRangeBillHandler{DB: db})
	router.Handle("/billrange", handler.RangeBillHandler{DB: db})
	router.Handle("/monthly", handler.MonthlyBillHandler{DB: db})
	router.Handle("/summary", handler.MonthlySummaryHandler{DB: db})
	router.Handle("/trend", handler.MonthlyTrendHandler{DB: db})
	router.Handle("/customertrend", handler.CustomerMonthlyTrendHandler{DB: db})
	router.Handle("/customerpersubtrend", handler.CustomerMonthlyTrendPerSubHandler{DB: db})
	router.Handle("/customermonthlypersub", handler.CustomerMonthlyCostPerSubHandler{DB: db})
	router.Handle("/customercostrangepersub", handler.CustomerRangeCostPerSubHandler{DB: db})
	router.Handle("/auth", handler.ADLoginHandler{Auth: auth})
	return router, nil

}

func startHTTPServer(c *config.Config, logger *ServiceLogger) (*http.Server, error) {

	srv, err := getServer(c, logger)
	if err != nil {
		return nil, err
	}
	router, err := getRouter(c, logger)
	if err != nil {
		return nil, err
	}
	srv.Handler = router

	go func() {
		if err := srv.ListenAndServeTLS(getPath(c.SSLCertPath), getPath(c.SSLKeyPath)); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Httpserver: ListenAndServe() error: %s", err), true)
			return
		}
		logger.Info("http server Shutting down", true)
	}()
	msg := fmt.Sprintf("http server Listening on %s", srv.Addr)
	logger.Info(msg, true)
	return srv, nil
}

func getIps(c *config.Config, logger *ServiceLogger) ([]net.IP, []net.IPNet, error) {
	iparr := c.APIClientIPs
	ipsarr := c.APIClientIPSubnets

	ips := make([]net.IP, len(iparr))
	ipns := make([]net.IPNet, len(ipsarr))
	var err error
	for i, p := range iparr {
		ips[i] = net.ParseIP(p)
		if ips[i] == nil {
			return nil, nil, fmt.Errorf("ip extraction error: %s", p)
		}
	}

	for i, p := range ipsarr {
		ipns[i], err = util.ParseIPNet(p)
		if err != nil {
			return nil, nil, err
		}
	}

	return ips, ipns, nil
}

func landing(resp http.ResponseWriter, req *http.Request) {

	resp.Write([]byte("Server UP"))
}

func getServer(c *config.Config, logger *ServiceLogger) (*http.Server, error) {
	addr := ":" + strconv.Itoa(c.Serverport)
	clientCAs := x509.NewCertPool()
	for _, cert := range c.ClientCAPool {
		certb, err := ioutil.ReadFile(getPath(cert))
		if err != nil {
			return nil, err
		}
		ok := clientCAs.AppendCertsFromPEM(certb)
		if !ok {
			logger.Warning("Coult not read client CA "+cert, true)
		}
	}
	tlsConf := &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  clientCAs,
	}

	return &http.Server{
		Addr:      addr,
		TLSConfig: tlsConf,
	}, nil

}

const (
	defConfigFilePath = "cspusage.json"
)
