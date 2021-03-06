package model

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/mr-tron/base58"
	uuid "github.com/satori/go.uuid"
)

var encoding = base32.NewEncoding("gxma8d698zf0cww77nhbg8jj1zv4sk8g")

// NewID is a globally unique identifier.
// It is generated by UUID version 4 and encoded with base58 with the [-] removed.
func NewID() string {

	uuidStr := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	uuidByte, _ := hex.DecodeString(uuidStr)
	uuidBase58 := base58.Encode(uuidByte)

	return uuidBase58
}

// NewRandomString randomly generates string with given length
func NewRandomString(length int) string {
	var b bytes.Buffer
	str := make([]byte, length+8)
	rand.Read(str)
	encoder := base32.NewEncoder(encoding, &b)
	encoder.Write(str)
	encoder.Close()
	b.Truncate(length) // removes the '==' padding
	return b.String()
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// GetMillis is a convenience method to get milliseconds since epoch.
func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// MapToJSON converts a map to a json string
func MapToJSON(objmap map[string]string) string {
	b, _ := json.Marshal(objmap)
	return string(b)
}

// MapFromJSON will decode the key/value pair map
func MapFromJSON(data io.Reader) map[string]string {
	decoder := json.NewDecoder(data)

	var objmap map[string]string
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]string)
	} else {
		return objmap
	}
}

func MapToQueryString(data map[string]interface{}) string {
	params := make([]string, 0)
	for k, v := range data {
		switch v := v.(type) {
		default:
			b, _ := json.Marshal(v)
			params = append(params, k+"="+string(b))
		case string:
			params = append(params, k+"="+data[k].(string))
		}
	}
	return "?" + strings.Join(params, "&")
}

func ArrayFromJSON(data io.Reader) []string {
	decoder := json.NewDecoder(data)

	var objmap []string
	if err := decoder.Decode(&objmap); err != nil {
		return make([]string, 0)
	} else {
		return objmap
	}
}

func CopyStringMap(originalMap map[string]string) map[string]string {
	copyMap := make(map[string]string)
	for k, v := range originalMap {
		copyMap[k] = v
	}
	return copyMap
}

// Copied from https://golang.org/src/net/dnsclient.go#L119
func IsDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
}

func IsLower(s string) bool {
	return strings.ToLower(s) == s
}
