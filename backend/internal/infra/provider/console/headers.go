package console

import (
	"net/http"
	"strings"

	"github.com/chenyme/grok2api/backend/internal/domain/account"
	infraegress "github.com/chenyme/grok2api/backend/internal/infra/egress"
	"github.com/chenyme/grok2api/backend/internal/infra/provider/browserheaders"
)

func applyBrowserHeaders(request *http.Request, token string, lease *infraegress.Lease, credential account.Credential) {
	accountKey := ""
	if lease != nil {
		accountKey = strings.TrimSpace(lease.AccountIdentity)
	}
	if accountKey == "" || accountKey == "shared" {
		accountKey = accountBrowserKey(credential.ID, credential.EgressIdentity)
	}
	identity := browserIdentityForAccount(accountKey)

	userAgent := ""
	if lease != nil {
		userAgent = strings.TrimSpace(lease.UserAgent)
	}
	if userAgent == "" {
		userAgent = strings.TrimSpace(identity.UserAgent)
	}
	if userAgent == "" {
		userAgent = infraegress.DefaultUserAgent
	}

	acceptLanguage := strings.TrimSpace(identity.AcceptLanguage)
	if acceptLanguage == "" {
		acceptLanguage = "zh-CN,zh;q=0.9,en;q=0.8"
	}
	acceptEncoding := strings.TrimSpace(identity.AcceptEncoding)
	if acceptEncoding == "" {
		acceptEncoding = "gzip, deflate, br, zstd"
	}

	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Encoding", acceptEncoding)
	request.Header.Set("Accept-Language", acceptLanguage)
	request.Header.Set("Cache-Control", "no-cache")
	cfCookies := ""
	if lease != nil {
		cfCookies = lease.CFCookies
	}
	request.Header.Set("Cookie", infraegress.BuildSSOCookie(token, cfCookies))
	request.Header.Set("Origin", "https://console.x.ai")
	request.Header.Set("Referer", "https://console.x.ai/")
	request.Header.Set("Sec-Fetch-Dest", "empty")
	request.Header.Set("Sec-Fetch-Mode", "cors")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Priority", "u=1, i")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("User-Agent", userAgent)
	applyChromiumClientHints(request.Header, userAgent)
}

// applyChromiumClientHints keeps the HTTP headers aligned with the Chromium
// TLS profile used by the Console transport. Non-Chromium User-Agents do not
// receive synthetic hints, avoiding contradictory browser fingerprints.
func applyChromiumClientHints(header http.Header, userAgent string) {
	browserheaders.ApplyChromiumClientHints(header, userAgent)
}
