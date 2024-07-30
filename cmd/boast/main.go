package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ciphermarco/BOAST/api"
	"github.com/ciphermarco/BOAST/config"
	"github.com/ciphermarco/BOAST/log"
	"github.com/ciphermarco/BOAST/receivers/dnsrcv"
	"github.com/ciphermarco/BOAST/receivers/httprcv"
	"github.com/ciphermarco/BOAST/storage"

	"github.com/BurntSushi/toml"
)

const program = "BOAST"
const version = "v1.0.0"
const author = "Marco Pereira (ciphermarco)"

var (
	prognver = fmt.Sprintf("%s %s", program, version)
	banner   = fmt.Sprintf("%s (by %s)\n", prognver, author)
	cfgPath  string
	logLevel int
	logPath  string
	dnsOnly  bool
	dnsTxt   string
	showVer  bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", banner)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "%s [OPTION...]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&cfgPath, "config", "boast.toml", "TOML configuration file")
	flag.IntVar(&logLevel, "log_level", 1, "Set the logging level (0=DEBUG|1=INFO)")
	flag.StringVar(&logPath, "log_file", "", "Path to log file")
	flag.BoolVar(&dnsOnly, "dns_only", false, "Run only the DNS receiver and its dependencies")
	flag.StringVar(&dnsTxt, "dns_txt", "", "DNS receiver's TXT record")
	flag.BoolVar(&showVer, "v", false, "Print program version and quit")
	flag.Parse()

	if showVer {
		fmt.Fprintf(os.Stderr, "%s", banner)
		os.Exit(0)
	}

	log.SetLevel(logLevel)
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		log.SetOutput(f)
	}
}

func main() {
	log.Info("Starting %s", prognver)

	tomlData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		log.Fatalln("Failed to read configuration:", err)
	}
	var cfg config.Config
	if err = toml.Unmarshal(tomlData, &cfg); err != nil {
		log.Fatalln("Failed to parse configuration:", err)
	}

	strg, err := storage.New(&storage.Config{
		TTL:             cfg.Strg.Expire.TTL.Value(),
		CheckInterval:   cfg.Strg.Expire.CheckInterval.Value(),
		MaxRestarts:     cfg.Strg.Expire.MaxRestarts,
		MaxEvents:       cfg.Strg.MaxEvents,
		MaxEventsByTest: cfg.Strg.MaxEventsByTest,
		MaxDumpSize:     cfg.Strg.MaxDumpSize.Value(),
		HMACKey:         cfg.Strg.HMACKey,
	})
	if err != nil {
		log.Fatalln("Failed to create storage:", err)
	}

	apiSrv := &api.Server{
		Host:        cfg.API.Host,
		Domain:      cfg.API.Domain,
		TLSPort:     cfg.API.TLSPort,
		TLSCertPath: cfg.API.TLSCertPath,
		TLSKeyPath:  cfg.API.TLSKeyPath,
		StatusPath:  cfg.API.Status.Path,
		Storage:     strg,
	}

	httpRcv := &httprcv.Receiver{
		Name:        "HTTP receiver",
		Host:        cfg.HTTPRcv.Host,
		Ports:       cfg.HTTPRcv.Ports,
		TLSPorts:    cfg.HTTPRcv.TLS.Ports,
		TLSCertPath: cfg.HTTPRcv.TLS.CertPath,
		TLSKeyPath:  cfg.HTTPRcv.TLS.KeyPath,
		IPHeader:    cfg.HTTPRcv.IPHeader,
		Storage:     strg,
	}

	txt := cfg.DNSRcv.Txt
	if dnsTxt != "" {
		txt = append(txt, dnsTxt)
	}
	dnsRcv := &dnsrcv.Receiver{
		Name:     "DNS receiver",
		Domain:   cfg.DNSRcv.Domain,
		Host:     cfg.DNSRcv.Host,
		Ports:    cfg.DNSRcv.Ports,
		PublicIP: cfg.DNSRcv.PublicIP,
		Txt:      txt,
		Storage:  strg,
	}

	errMain := make(chan error, 1)

	go strg.StartExpire(errMain)
	go dnsRcv.ListenAndServe(errMain)

	if !dnsOnly {
		go apiSrv.ListenAndServe(errMain)
		go httpRcv.ListenAndServe(errMain)
	}

	if exitErr := <-errMain; exitErr != nil {
		log.Info("Fatal error")
		log.Debug("Error: %v", exitErr)
		os.Exit(1)
	}
}
