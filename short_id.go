package shortid

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	redis "github.com/go-redis/redis"
)

const (
	all           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	timeFormat    = "060102"
	numberFormat  = "%08d"
	counterPrefix = "counter:"
	expireTime    = time.Hour * 48
)

var (
	_client     *redis.Client
	_serverHash string
	_chars      = [62]rune{}
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	_serverHash = sha256hash(hostname)[0:2]
	for i, char := range all {
		_chars[i] = char
	}
}

type Options struct {
	Number        int
	StartWithYear bool
	EndWithHost   bool
}

func Generate(opt Options) string {
	data, _ := generateRandomBytes(opt.Number)

	var buffer bytes.Buffer
	if opt.StartWithYear {
		year := time.Now().UTC().Format("06")
		buffer.WriteString(year)
	}

	for _, b := range data {
		pick := b % 61
		buffer.WriteRune(_chars[pick])
	}

	if opt.EndWithHost {
		buffer.WriteString(_serverHash)
	}

	return buffer.String()
}

func sha256hash(text string) string {
	rawBytes := []byte(text)
	h := sha256.Sum256(rawBytes)
	return base64.StdEncoding.EncodeToString(h[:])
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

//RedisConfig the configuration of redis
type RedisConfig redis.Options

//SetRedis setting the counter redis
func SetRedis(opt RedisConfig) {
	rOpt := redis.Options(opt)
	_client = redis.NewClient(&rOpt)
}

//GetCounter get the date serial number return int 64 (Example: 170817000000001)
// It will return an error if the redis error
func GetCounter(prefixStr string) (int64, error) {
	timeStr := time.Now().UTC().Format(timeFormat)
	rKey := counterPrefix + prefixStr + ":" + timeStr
	result, err := _client.Incr(rKey).Result()
	if err != nil {
		return -1, err
	}
	if result == 1 {
		_client.Expire(rKey, expireTime)
	}
	timeStr += numberFormat
	idStr := fmt.Sprintf(timeStr, result)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}
	return id, nil
}
