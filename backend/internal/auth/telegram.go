package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

var ErrInvalidInitData = errors.New("invalid telegram init data")

func ValidateInitData(initData string, botToken string) (url.Values, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init data: %w", err)
	}

	receivedHash := values.Get("hash")
	if receivedHash == "" {
		return nil, ErrInvalidInitData
	}
	values.Del("hash")

	// собрать data-check-string: все пары key=value, отсортированные по ключу, через \n
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+values.Get(k))
	}
	dataCheckString := strings.Join(pairs, "\n")

	// секретный ключ = HMAC-SHA256("WebAppData", botToken)
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))

	h := hmac.New(sha256.New, secretKey.Sum(nil))
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	if computedHash != receivedHash {
		return nil, ErrInvalidInitData
	}

	return values, nil
}