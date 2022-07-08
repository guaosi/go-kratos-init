package process

import (
	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	charset          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lockStringLength = 6
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// 数字，大小写字母
func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

// "1,2, 3" => []string{"1", "2", "3"}
func SplitTrimSpace(str, sep string) []string {
	var result []string

	for _, s := range strings.Split(str, sep) {
		result = append(result, strings.TrimSpace(s))
	}

	return result
}

// 简化uuid, 20个字符，64位
func SimpleUUID() string {
	return xid.New().String()
}

// 标准uuid，36个字符，128位
func UUID() string {
	return uuid.NewV4().String()
}

// 1000 => 1,000
func Comma(v int64) string {
	sign := ""

	// Min int64 can't be negated to a usable value, so it has to be special cased.
	if v == math.MinInt64 {
		return "-9,223,372,036,854,775,808"
	}

	if v < 0 {
		sign = "-"
		v = 0 - v
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for v > 999 {
		parts[j] = strconv.FormatInt(v%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		v = v / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(v))
	return sign + strings.Join(parts[j:], ",")
}

func CommaUint(v uint) string {
	return Comma(int64(v))
}

// 分布式锁随机值
func GetLockRandStringID() string {
	str := UUID()
	return str[len(str)-lockStringLength:]
}

func TimeTransFormString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetNowTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// convert []string to []int
func ConvertToInts(s []string) ([]int, error) {
	data := make([]int, 0, len(s))

	for _, e := range s {
		d, err := strconv.Atoi(e)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}

	return data, nil
}
