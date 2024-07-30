package config

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Config represents BOAST's configuration.
// It contains the structs for each configuration section and is used to unmarshal the
// TOML configuration file.
type Config struct {
	API     APIConfig     `toml:"api"`
	HTTPRcv HTTPRcvConfig `toml:"http_receiver"`
	DNSRcv  DNSRcvConfig  `toml:"dns_receiver"`
	Strg    StorageConfig `toml:"storage"`
}

// APIConfig represents the web API configuration.
type APIConfig struct {
	Host        string          `toml:"host"`
	Domain      string          `toml:"domain"`
	TLSPort     int             `toml:"tls_port"`
	TLSCertPath string          `toml:"tls_cert"`
	TLSKeyPath  string          `toml:"tls_key"`
	Status      APIStatusConfig `toml:"status"`
}

// APIStatusConfig represents the web API configuration specific to the status page.
type APIStatusConfig struct {
	Path string `toml:"url_path"`
}

// HTTPRcvConfig represents the HTTP protocol receiver configuration.
type HTTPRcvConfig struct {
	Host     string           `toml:"host"`
	Ports    []int            `toml:"ports"`
	TLS      HTTPRcvConfigTLS `toml:"tls"`
	IPHeader string           `toml:"real_ip_header"`
}

// HTTPRcvConfigTLS represents the HTTP protocol receiver configuration specific to its
// TLS functionalities.
type HTTPRcvConfigTLS struct {
	Ports    []int  `toml:"ports"`
	CertPath string `toml:"cert"`
	KeyPath  string `toml:"key"`
}

// DNSRcvConfig represents the DNS protocol receiver configuration.
type DNSRcvConfig struct {
	Domain   string   `toml:"domain"`
	Host     string   `toml:"host"`
	Ports    []int    `toml:"ports"`
	PublicIP string   `toml:"public_ip"`
	Txt      []string `toml:"txt"`
}

// StorageConfig represents the storage configuration.
type StorageConfig struct {
	MaxEvents       int          `toml:"max_events"`
	MaxEventsByTest int          `toml:"max_events_by_test"`
	MaxDumpSize     byteSize     `toml:"max_dump_size"`
	HMACKey         hmacKey      `toml:"hmac_key"`
	Expire          ExpireConfig `toml:"expire"`
}

// ExpireConfig represents the storage configurations specific to its expiration feature.
type ExpireConfig struct {
	TTL           duration `toml:"ttl"`
	CheckInterval duration `toml:"check_interval"`
	MaxRestarts   int      `toml:"max_restarts"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(txt []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(txt))
	return err
}

func (d *duration) Value() time.Duration {
	return d.Duration
}

type hmacKey []byte

func (k *hmacKey) UnmarshalText(txt []byte) error {
	maxKeySize := 64
	if len(txt) > maxKeySize {
		return errors.New("hmac_key must be between 0 and 64 bytes long")
	}
	*k = txt
	return nil
}

// byteSize is the type representing the value in bytes for a given unit string (e.g. "80KB").
// It's an int since only whole bytes will be considered and it's not realistic that
// somebody will need EiB here (which overflows with int).
type byteSize int

func (b *byteSize) UnmarshalText(txt []byte) error {
	bs, err := parseByteSize(string(txt))
	if err != nil {
		return err
	}
	*b = bs
	return nil
}

func (b *byteSize) Value() int {
	return int(*b)
}

const (
	// B represents 1 byte.
	B byteSize = 1

	// KiB is the number of bytes in 1 kibibyte.
	KiB = 1 << (10 * iota)
	// MiB is the number of bytes in 1 mebibyte.
	MiB
	// GiB is the number of bytes in 1 gibibyte.
	GiB
	// TiB is the number of bytes in 1 tebibyte.
	TiB
	// PiB is the number of bytes in 1 pebibyte.
	PiB
	// EiB overflows with int (reasons to use int are above byteSize's declaration).
	// This means you cannot use your black hole computer's full capacity.

	// KB is the number of bytes in 1 kilobyte.
	KB byteSize = 1e3
	// MB is the number of bytes in 1 megabyte.
	MB byteSize = 1e6
	// GB is the number of bytes in 1 gigabyte.
	GB byteSize = 1e9
	// TB is the number of bytes in 1 terabyte.
	TB byteSize = 1e12
	// PB is the number of bytes in 1 petabyte.
	PB byteSize = 1e15
	// EB is the number of bytes in 1 exabyte.
	EB byteSize = 1e18
)

var unitToByteSize = map[string]byteSize{
	"B": B,

	"KIB": KiB,
	"MIB": MiB,
	"GIB": GiB,
	"TIB": TiB,
	"PIB": PiB,

	"KB": KB,
	"MB": MB,
	"GB": GB,
	"TB": TB,
	"PB": PB,
	"EB": EB,
}

func parseByteSize(s string) (byteSize, error) {
	s = strings.TrimSpace(s)
	var ss []string
	for i, c := range s {
		if !unicode.IsDigit(c) && c != '.' {
			ss = append(ss, strings.TrimSpace(string(s[:i])))
			ss = append(ss, strings.TrimSpace(string(s[i:])))
			break
		}
	}

	if len(ss) != 2 {
		return 0, errors.New("wrong format")
	}

	unit, exists := unitToByteSize[strings.ToUpper(ss[1])]
	if !exists {
		return 0, errors.New("unrecognised size suffix " + ss[1])
	}

	sn := ss[0]
	var bs byteSize
	if strings.Contains(sn, ".") {
		n, err := strconv.ParseFloat(sn, 64)
		if err != nil {
			return 0, err
		}
		bs = byteSize(n * float64(unit))
	} else {
		n, err := strconv.Atoi(sn)
		if err != nil {
			return 0, err
		}
		bs = byteSize(n * int(unit))
	}

	return bs, nil
}
