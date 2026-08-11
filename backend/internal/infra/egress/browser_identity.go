package egress

import (
	"crypto/rand"
	"math/big"
	"strings"
	"sync"
)

// BrowserProfile is a sticky browser surface bound to one account identity.
type BrowserProfile struct {
	UserAgent      string
	AcceptLanguage string
	AcceptEncoding string
}

type stickyBrowserPlatform struct {
	userAgent string
	language  string
}

var stickyBrowserPlatforms = []stickyBrowserPlatform{
	{
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
		language:  "en-US,en;q=0.9",
	},
	{
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
		language:  "en-US,en;q=0.9",
	},
	{
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36",
		language:  "en-GB,en;q=0.9",
	},
	{
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36",
		language:  "zh-CN,zh;q=0.9,en;q=0.8",
	},
	{
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
		language:  "zh-CN,zh;q=0.9,en;q=0.8",
	},
	{
		userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36",
		language:  "ja-JP,ja;q=0.9,en-US;q=0.8,en;q=0.7",
	},
}

var stickyBrowserEncodings = []string{
	"gzip, deflate, br, zstd",
	"gzip, deflate, br",
	"gzip, deflate",
}

var stickyBrowserByAccount sync.Map // map[string]BrowserProfile

// StickyBrowserProfile returns a stable browser surface for one account identity.
// Different accounts receive different profiles; the same account keeps the same one.
func StickyBrowserProfile(accountKey string) BrowserProfile {
	accountKey = strings.TrimSpace(accountKey)
	if accountKey == "" {
		return BrowserProfile{
			UserAgent:      DefaultUserAgent,
			AcceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			AcceptEncoding: "gzip, deflate, br, zstd",
		}
	}
	if cached, ok := stickyBrowserByAccount.Load(accountKey); ok {
		if profile, okProfile := cached.(BrowserProfile); okProfile {
			return profile
		}
	}
	profile := newStickyBrowserProfile()
	actual, _ := stickyBrowserByAccount.LoadOrStore(accountKey, profile)
	if stored, ok := actual.(BrowserProfile); ok {
		return stored
	}
	return profile
}

// StickyBrowserUserAgent is a convenience wrapper for transport selection.
func StickyBrowserUserAgent(accountKey string) string {
	return StickyBrowserProfile(accountKey).UserAgent
}

func newStickyBrowserProfile() BrowserProfile {
	platform := pickStickyBrowserPlatform()
	return BrowserProfile{
		UserAgent:      platform.userAgent,
		AcceptLanguage: platform.language,
		AcceptEncoding: pickStickyBrowserString(stickyBrowserEncodings),
	}
}

func pickStickyBrowserPlatform() stickyBrowserPlatform {
	if len(stickyBrowserPlatforms) == 0 {
		return stickyBrowserPlatform{userAgent: DefaultUserAgent, language: "en-US,en;q=0.9"}
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(stickyBrowserPlatforms))))
	if err != nil {
		return stickyBrowserPlatforms[0]
	}
	return stickyBrowserPlatforms[int(n.Int64())]
}

func pickStickyBrowserString(options []string) string {
	if len(options) == 0 {
		return ""
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	if err != nil {
		return options[0]
	}
	return options[int(n.Int64())]
}

// InvalidateStickyBrowserProfile drops one account's cached browser surface.
func InvalidateStickyBrowserProfile(accountKey string) {
	accountKey = strings.TrimSpace(accountKey)
	if accountKey == "" {
		return
	}
	stickyBrowserByAccount.Delete(accountKey)
}
