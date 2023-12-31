package unique

// Package unique 产生一个唯一字符串

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"share.ac.cn/common/autoinc"
)

var stringInst, numberInst, dateInst *Unique

// Unique 基于时间戳的唯一不定长字符串
//
// NOTE: 算法是基于系统时间的。所以必须得保证时间上正确的，否则可能会造成非唯一的情况。
// NOTE: 产生的数据有一定的顺序性。
//
// Unique 由两部分组成：
// 前缀是由一个相对稳定的字符串，与时间相关联；
// 后缀是一个自增的数值。
//
// 每次刷新前缀之后，都会重置后缀的计数器，从头开始。
// 刷新时间和计数器的步长都是一个随机数。
type Unique struct {
	random *rand.Rand

	// 数据转换成字符串所采用的进制。
	formatBase int

	// 前缀部分的内容。
	//
	// 根据 prefixFormat 是否存在，会呈现不同的内容：
	// 如果 prefixFormat 为空，prefix 为一个时间戳的整数值，
	// 按一定的进制进行转换之后的值；否则是按 prefixFormat
	// 进行格式化的时间数据。
	prefix       string
	prefixFormat string

	timer    *time.Timer
	duration time.Duration

	step int64
	ai   *autoinc.AutoInc

	// 用保证 prefix 和 ai 的一致性。
	resetLocker sync.RWMutex
}

// String 返回以字符串形式表示的 Unique 实例
//
// 格式为：p4k5f81
//
// NOTE: 多次调用，返回的是同一个实例。
func String() *Unique {
	if stringInst == nil {
		stringInst = NewString()
	}

	return stringInst
}

// NewString 声明以字符串形式表示的 Unique 实例
//
// 格式为：p4k5f81
//
// 与 String 的不同在于，每次调用 NewString 都返回新的实例，而 String 则是返回相同实例。
func NewString() *Unique {
	return New(time.Now().Unix(), 1, time.Hour, "", 36)
}

// Number 返回数字形式表示的 Unique 实例
//
// 格式为：15193130121
//
// NOTE: 多次调用，返回的是同一个实例。
func Number() *Unique {
	if numberInst == nil {
		numberInst = NewNumber()
	}

	return numberInst
}

// NewNumber 声明以数字形式表示的 Unique 实例
//
// 格式为：15193130121
//
// 与 Number 的不同在于，每次调用 NewNumber 都返回新的实例，而 Number 则是返回相同实例。
func NewNumber() *Unique {
	return New(time.Now().Unix(), 1, time.Hour, "", 10)
}

// Date 返回以日期形式表示的 Unique 实例
//
// 格式为：20180222232332-1
//
// NOTE: 多次调用，返回的是同一个实例。
func Date() *Unique {
	if dateInst == nil {
		dateInst = NewDate()
	}

	return dateInst
}

// NewDate 声明以日期形式表示的 Unique 实例
//
// 格式为：20180222232332-1
//
// 与 Date 的不同在于，每次调用 NewDate 都返回新的实例，而 Date 则是返回相同实例。
func NewDate() *Unique {
	return New(time.Now().Unix(), 1, time.Hour, "20060102150405-", 10)
}

// New 声明一个新的 Unique。
//
// seed 随机种子；
// step 计数器的步长，需大于 0；
// duration 计数器的重置时间，不能小于 1*time.Second；
// prefixFormat 格式化 prefix 的方式，若指定，则格式化为时间，否则将时间戳转换为数值；
// base 数值转换成字符串时，所采用的进制，可以是 [2,36] 之间的值。
func New(seed, step int64, duration time.Duration, prefixFormat string, base int) *Unique {
	if step <= 0 {
		panic("无效的参数 step")
	}

	if duration < time.Second {
		panic("无效的参数 duration，不能小于 1 秒")
	}

	if prefixFormat != "" && !isValidDateFormat(prefixFormat) {
		panic("无效的 prefixFormat 参数")
	}

	if base < 2 || base > 36 {
		panic("无效的 base 值，只能介于 [2,36] 之间")
	}

	u := &Unique{
		random:       rand.New(rand.NewSource(seed)),
		formatBase:   base,
		duration:     duration,
		prefixFormat: prefixFormat,
		step:         step,
	}

	u.reset()

	return u
}

func isValidDateFormat(format string) bool {
	return strings.Contains(format, "2006") &&
		strings.Contains(format, "01") &&
		strings.Contains(format, "02") &&
		strings.Contains(format, "15") &&
		strings.Contains(format, "04") &&
		strings.Contains(format, "05")
}

// 重置时间戳和计数器
func (u *Unique) reset() {
	u.resetLocker.Lock()
	defer u.resetLocker.Unlock()

	if u.prefixFormat != "" {
		u.prefix = time.Now().Format(u.prefixFormat)
	} else {
		u.prefix = strconv.FormatInt(time.Now().Unix(), u.formatBase)
	}

	if u.ai != nil {
		u.ai.Stop()
	}
	u.ai = autoinc.New(1, u.step, 1000)

	if u.timer != nil {
		u.timer.Stop()
	}

	u.timer = time.AfterFunc(u.duration, u.reset)
}

// String 返回一个唯一的字符串
func (u *Unique) String() string {
	u.resetLocker.RLock()
	p := u.prefix
	id, ok := u.ai.ID()
	u.resetLocker.RUnlock()

	for !ok {
		u.reset() // NOTE: reset 包含对 resetLocker 的操作

		u.resetLocker.RLock()
		p = u.prefix
		id, ok = u.ai.ID()
		u.resetLocker.RUnlock()
	}

	return p + strconv.FormatInt(id, u.formatBase)
}

// Bytes 返回 String() 的 []byte 格式
//
// 在多次出错之后，可能会触发 panic
func (u *Unique) Bytes() []byte { return []byte(u.String()) }
