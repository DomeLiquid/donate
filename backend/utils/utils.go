package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

func newUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}

func GenUuidFromStrings(uuids ...string) string {
	if len(uuids) == 0 {
		uuids = append(uuids, "00000000-0000-0000-0000-000000000000")
	}

	// Sort the UUIDs to ensure consistent ordering
	sortedUUIDs := make([]string, len(uuids))
	copy(sortedUUIDs, uuids)
	sort.Strings(sortedUUIDs)

	// Concatenate all sorted UUIDs
	concatenatedUUIDs := strings.Join(sortedUUIDs, "")

	return uuidHash([]byte(concatenatedUUIDs))
}

func uuidHash(b []byte) string {
	h := md5.New()

	h.Write(b)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}

func RandomPin() string {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	c := binary.LittleEndian.Uint64(b[:]) % 1000000
	if c < 100000 {
		c = 100000 + c
	}

	return strconv.FormatUint(c, 10)
}

func RandomTraceID() string {
	return newUUID()
}

type PaymentParams struct {
	UUID     string
	Asset    string
	Amount   string
	Memo     string
	Trace    string
	ReturnTo string
}

func BuildMixinOneSafePaymentURI(params PaymentParams) string {
	address := params.UUID
	if params.UUID != "" && uuid.FromStringOrNil(params.UUID) != uuid.Nil {
		address = params.UUID
	}

	baseURL := fmt.Sprintf("https://mixin.one/pay/%s", address)

	query := url.Values{}
	query.Set("asset", params.Asset)
	query.Set("amount", params.Amount)
	query.Set("memo", params.Memo)

	if params.Trace != "" {
		query.Set("trace", params.Trace)
	} else {
		query.Set("trace", newUUID())
	}

	if params.ReturnTo != "" {
		query.Set("return_to", url.QueryEscape(params.ReturnTo))
	}

	return fmt.Sprintf("%s?%s", baseURL, query.Encode())
}
