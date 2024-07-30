package config_test

import (
	"bytes"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/ciphermarco/BOAST/config"
)

var data = []byte(`
[api]
  host = "0.0.0.0"
  tls_port = 2096
  tls_cert = "/path/to/tls/server.crt"
  tls_key = "/path/to/tls/server.key"

  [api.status]
    url_path = "rzaedgmqloivvw7v3lamu3tzvi"

[http_receiver]
  host = "0.0.0.0"
  ports = [80, 8080]
  real_ip_header = "X-Real-IP"

  [http_receiver.tls]
    ports = [443, 8443]
    cert = "/path/to/tls/server.crt"
    key = "/path/to/tls/server.key"

[dns_receiver]
  domain = "example.com"
  host = "0.0.0.0"
  ports = [53, 5353]
  public_ip = "203.0.113.77"

[storage]
  max_events = 1_000_000
  max_events_by_test = 100
  max_dump_size = "80KB"
  hmac_key = "TJkhXnMqSqOaYDiTw7HsfQ=="

  [storage.expire]
    ttl = "24h"
    check_interval = "1h"
    max_restarts = 100
`)

func TestConfigParse(t *testing.T) {
	var cfg config.Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("config parsing error: %v (want) != %v (got)",
			nil, err)
	}

	// API
	wantAPIHost := "0.0.0.0"
	gotAPIHost := cfg.API.Host
	if wantAPIHost != gotAPIHost {
		t.Errorf("wrong API host: %v (want) != %v (got)",
			wantAPIHost, gotAPIHost)
	}

	wantAPITLSPort := 2096
	gotAPITLSPort := cfg.API.TLSPort
	if wantAPITLSPort != gotAPITLSPort {
		t.Errorf("wrong API TLS Port: %v (want) != %v (got)",
			wantAPITLSPort, gotAPITLSPort)
	}

	wantAPITLSCertPath := "/path/to/tls/server.crt"
	gotAPITLSCertPath := cfg.API.TLSCertPath
	if wantAPITLSCertPath != gotAPITLSCertPath {
		t.Errorf("wrong API TLS cert path: %v (want) != %v (got)",
			wantAPITLSCertPath, gotAPITLSCertPath)
	}

	wantAPITLSKeyPath := "/path/to/tls/server.key"
	gotAPITLSKeyPath := cfg.API.TLSKeyPath
	if wantAPITLSKeyPath != gotAPITLSKeyPath {
		t.Errorf("wrong API TLS key path: %v (want) != %v (got)",
			wantAPITLSKeyPath, gotAPITLSKeyPath)
	}

	wantAPIStatusPath := "rzaedgmqloivvw7v3lamu3tzvi"
	gotAPIStatusPath := cfg.API.Status.Path
	if wantAPIStatusPath != gotAPIStatusPath {
		t.Errorf("wrong API status URL path: %v (want) != %v (got)",
			wantAPIStatusPath, gotAPIStatusPath)
	}

	// Storage
	wantStrgMaxEvents := 1_000_000
	gotStrgMaxEvents := cfg.Strg.MaxEvents
	if wantStrgMaxEvents != gotStrgMaxEvents {
		t.Errorf("wrong Storage max events: %v (want) != %v (got)",
			wantStrgMaxEvents, gotStrgMaxEvents)
	}

	wantStrgMaxEventsByTest := 100
	gotStrgMaxEventsByTest := cfg.Strg.MaxEventsByTest
	if wantStrgMaxEventsByTest != gotStrgMaxEventsByTest {
		t.Errorf("wrong Storage max events by test: %v (want) != %v (got)",
			wantStrgMaxEventsByTest, gotStrgMaxEventsByTest)
	}

	wantStrgMaxDumpSize := int(80 * 1e3) // "80KB"
	gotStrgMaxDumpSize := cfg.Strg.MaxDumpSize.Value()
	if wantStrgMaxDumpSize != gotStrgMaxDumpSize {
		t.Errorf("wrong Storage max dump size: %v (want) != %v (got)",
			wantStrgMaxDumpSize, gotStrgMaxDumpSize)
	}

	wantStrgHMACKey := []byte("TJkhXnMqSqOaYDiTw7HsfQ==")
	gotStrgHMACKey := cfg.Strg.HMACKey
	if !bytes.Equal(wantStrgHMACKey, gotStrgHMACKey) {
		t.Errorf("wrong Storage HMAC key: %v (want) != %v (got)",
			wantStrgHMACKey, gotStrgHMACKey)
	}

	wantStrgTTL := time.Duration(24 * time.Hour)
	gotStrgTTL := cfg.Strg.Expire.TTL.Value()
	if wantStrgTTL != gotStrgTTL {
		t.Errorf("wrong Storage TTL: %v (want) != %v (got)",
			wantStrgTTL, gotStrgTTL)
	}

	wantStrgCheckInterval := time.Duration(1 * time.Hour)
	gotStrgCheckInterval := cfg.Strg.Expire.CheckInterval.Value()
	if wantStrgCheckInterval != gotStrgCheckInterval {
		t.Errorf("wrong Storage check interval: %v (want) != %v (got)",
			wantStrgCheckInterval, gotStrgCheckInterval)
	}

	wantStrgMaxRestarts := 100
	gotStrgMaxRestarts := cfg.Strg.Expire.MaxRestarts
	if wantStrgMaxRestarts != gotStrgMaxRestarts {
		t.Errorf("wrong Storage max restarts: %v (want) != %v (got)",
			wantStrgMaxRestarts, gotStrgMaxRestarts)
	}

	// HTTP Receiver
	wantHTTPRcvHost := "0.0.0.0"
	gotHTTPRcvHost := cfg.HTTPRcv.Host
	if wantHTTPRcvHost != gotHTTPRcvHost {
		t.Errorf("wrong HTTP receiver host: %v (want) != %v (got)",
			wantHTTPRcvHost, gotHTTPRcvHost)
	}

	wantHTTPRcvPorts := []int{80, 8080}
	gotHTTPRcvPorts := cfg.HTTPRcv.Ports
	if !reflect.DeepEqual(wantHTTPRcvPorts, gotHTTPRcvPorts) {
		t.Errorf("wrong HTTP receiver ports: %v (want) != %v (got)",
			wantHTTPRcvPorts, gotHTTPRcvPorts)
	}

	wantHTTPRcvIPHeader := "X-Real-IP"
	gotHTTPRcvIPHeader := cfg.HTTPRcv.IPHeader
	if wantHTTPRcvIPHeader != gotHTTPRcvIPHeader {
		t.Errorf("wrong HTTP receiver IP header: %v (want) != %v (got)",
			wantHTTPRcvIPHeader, gotHTTPRcvIPHeader)
	}

	wantHTTPRcvTLSPorts := []int{443, 8443}
	gotHTTPRcvTLSPorts := cfg.HTTPRcv.TLS.Ports
	if !reflect.DeepEqual(wantHTTPRcvTLSPorts, gotHTTPRcvTLSPorts) {
		t.Errorf("wrong HTTP receiver TLS ports: %v (want) != %v (got)",
			wantHTTPRcvTLSPorts, gotHTTPRcvTLSPorts)
	}

	wantHTTPRcvTLSCertPath := "/path/to/tls/server.crt"
	gotHTTPRcvTLSCertPath := cfg.HTTPRcv.TLS.CertPath
	if wantHTTPRcvTLSCertPath != gotHTTPRcvTLSCertPath {
		t.Errorf("wrong HTTP receiver TLS cert path: %v (want) != %v (got)",
			wantHTTPRcvTLSCertPath, gotHTTPRcvTLSCertPath)
	}

	wantHTTPRcvTLSKeyPath := "/path/to/tls/server.key"
	gotHTTPRcvTLSKeyPath := cfg.HTTPRcv.TLS.KeyPath
	if wantHTTPRcvTLSKeyPath != gotHTTPRcvTLSKeyPath {
		t.Errorf("wrong HTTP receiver TLS key path: %v (want) != %v (got)",
			wantHTTPRcvTLSKeyPath, gotHTTPRcvTLSKeyPath)
	}

	// DNS Receiver
	wantDNSRcvDomain := "example.com"
	gotDNSRcvDomain := cfg.DNSRcv.Domain
	if wantDNSRcvDomain != gotDNSRcvDomain {
		t.Errorf("wrong DNS receiver domain: %v (want) != %v (got)",
			wantDNSRcvDomain, gotDNSRcvDomain)
	}

	wantDNSRcvHost := "0.0.0.0"
	gotDNSRcvHost := cfg.DNSRcv.Host
	if wantDNSRcvHost != gotDNSRcvHost {
		t.Errorf("wrong DNS receiver host: %v (want) != %v (got)",
			wantDNSRcvHost, gotDNSRcvHost)
	}

	wantDNSRcvPorts := []int{53, 5353}
	gotDNSRcvPorts := cfg.DNSRcv.Ports
	if !reflect.DeepEqual(wantDNSRcvPorts, gotDNSRcvPorts) {
		t.Errorf("wrong DNS receiver ports: %v (want) != %v (got)",
			wantDNSRcvPorts, gotDNSRcvPorts)
	}

	wantDNSRcvPublicIP := "203.0.113.77"
	gotDNSRcvPublicIP := cfg.DNSRcv.PublicIP
	if wantDNSRcvPublicIP != gotDNSRcvPublicIP {
		t.Errorf("wrong DNS receiver public IP: %v (want) != %v (got)",
			wantDNSRcvPublicIP, gotDNSRcvPublicIP)
	}
}

func TestStorageMaxDumpFloatRounding(t *testing.T) {
	var maxDumpSize = []byte(
		`[storage]
		   # an invalid float syntax will make ParseFloat fail so the error
		   # code path can be reached
                   max_dump_size = "80.5555555KB"`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(maxDumpSize, &cfg); err != nil {
		t.Fatalf("unexpected error: %v (want) != %v (got)",
			nil, err)
	}

	wantStrgMaxDumpSize := 80555 // from 80555.5555
	gotStrgMaxDumpSize := cfg.Strg.MaxDumpSize.Value()
	if wantStrgMaxDumpSize != gotStrgMaxDumpSize {
		t.Errorf("wrong max_dump_size: %v (want) != %v (got)",
			wantStrgMaxDumpSize, gotStrgMaxDumpSize)
	}
}

func TestStorageHMACKeyParseError(t *testing.T) {
	var tooLongHMACKey = []byte(
		`[storage]
                   hmac_key = "UtFdm4qQa56yZEfwWEWf1NG/IJKzUya6jYtWCWKqjAclUaiEI5hXh9LrBfJrWEkmM/dXnvxiDgfHeD+EjkRCEpY="`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(tooLongHMACKey, &cfg); err == nil {
		t.Errorf("parsing hmac_key did not fail: error (want) != %v (got)",
			err)
	}
}

func TestStorageMaxDumpSizeFormatError(t *testing.T) {
	var badMaxDumpSize = []byte(
		`[storage]
                   max_dump_size = "80"`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(badMaxDumpSize, &cfg); err == nil {
		t.Errorf("parsing max_dump_size did not fail: error (want) != %v (got)",
			err)
	}
}

func TestStorageMaxDumpSizeSuffixError(t *testing.T) {
	var badMaxDumpSize = []byte(
		`[storage]
                   max_dump_size = "80KBB"`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(badMaxDumpSize, &cfg); err == nil {
		t.Errorf("parsing max_dump_size did not fail: error (want) != %v (got)",
			err)
	}
}

func TestStorageMaxDumpParseFloatError(t *testing.T) {
	var badMaxDumpSize = []byte(
		`[storage]
		   # an invalid float syntax will make ParseFloat fail so the error
		   # code path can be reached
                   max_dump_size = "1.0000000000000001110223024625156540423631668090820312500...001KB"`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(badMaxDumpSize, &cfg); err == nil {
		t.Errorf("parsing max_dump_size did not fail: error (want) != %v (got)",
			err)
	}
}

func TestStorageMaxDumpSizeAtoiError(t *testing.T) {
	var badMaxDumpSize = []byte(
		`[storage]
		   # a too big int will make Atoi fail so the error code path
		   # can be reached
                   max_dump_size = "99999999999999999999999999999999999999KB"`,
	)
	var cfg config.Config
	if err := toml.Unmarshal(badMaxDumpSize, &cfg); err == nil {
		t.Errorf("parsing max_dump_size did not fail: error (want) != %v (got)",
			err)
	}
}

func TestParseByteSize(t *testing.T) {
	wantB := 1
	gotB, err := config.ParseByteSize("1B")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantB != int(gotB) {
		t.Errorf("wrong B byte size: %v (want) != %v (got)",
			wantB, gotB)
	}

	wantKiB := 1024
	gotKiB, err := config.ParseByteSize("1KiB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantKiB != int(gotKiB) {
		t.Errorf("wrong KiB byte size: %v (want) != %v (got)",
			wantKiB, gotKiB)
	}

	wantMiB := int(math.Pow(1024, 2))
	gotMiB, err := config.ParseByteSize("1MiB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantMiB != int(gotMiB) {
		t.Errorf("wrong MiB byte size: %v (want) != %v (got)",
			wantMiB, gotMiB)
	}

	wantGiB := int(math.Pow(1024, 3))
	gotGiB, err := config.ParseByteSize("1GiB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantGiB != int(gotGiB) {
		t.Errorf("wrong GiB byte size: %v (want) != %v (got)",
			wantGiB, gotGiB)
	}

	wantTiB := int(math.Pow(1024, 4))
	gotTiB, err := config.ParseByteSize("1TiB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantTiB != int(gotTiB) {
		t.Errorf("wrong TiB byte size: %v (want) != %v (got)",
			wantTiB, gotTiB)
	}

	wantPiB := int(math.Pow(1024, 5))
	gotPiB, err := config.ParseByteSize("1PiB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantPiB != int(gotPiB) {
		t.Errorf("wrong PiB byte size: %v (want) != %v (got)",
			wantPiB, gotPiB)
	}

	wantKB := 1000
	gotKB, err := config.ParseByteSize("1KB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantKB != int(gotKB) {
		t.Errorf("wrong KB byte size: %v (want) != %v (got)",
			wantKB, gotKB)
	}

	wantMB := int(math.Pow(1000, 2))
	gotMB, err := config.ParseByteSize("1MB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantMB != int(gotMB) {
		t.Errorf("wrong MB byte size: %v (want) != %v (got)",
			wantMB, gotMB)
	}

	wantGB := int(math.Pow(1000, 3))
	gotGB, err := config.ParseByteSize("1GB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantGB != int(gotGB) {
		t.Errorf("wrong GB byte size: %v (want) != %v (got)",
			wantGB, gotGB)
	}

	wantTB := int(math.Pow(1000, 4))
	gotTB, err := config.ParseByteSize("1TB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantTB != int(gotTB) {
		t.Errorf("wrong TB byte size: %v (want) != %v (got)",
			wantTB, gotTB)
	}

	wantPB := int(math.Pow(1000, 5))
	gotPB, err := config.ParseByteSize("1PB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantPB != int(gotPB) {
		t.Errorf("wrong PB byte size: %v (want) != %v (got)",
			wantPB, gotPB)
	}

	wantEB := int(math.Pow(1000, 6))
	gotEB, err := config.ParseByteSize("1EB")
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantEB != int(gotEB) {
		t.Errorf("wrong EB byte size: %v (want) != %v (got)",
			wantEB, gotEB)
	}
}
