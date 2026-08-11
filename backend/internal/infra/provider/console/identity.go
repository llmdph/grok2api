package console

import (
	"fmt"
	"strings"

	infraegress "github.com/chenyme/grok2api/backend/internal/infra/egress"
)

type browserIdentity struct {
	UserAgent      string
	AcceptLanguage string
	AcceptEncoding string
}

func browserIdentityForAccount(accountKey string) browserIdentity {
	profile := infraegress.StickyBrowserProfile(accountKey)
	return browserIdentity{
		UserAgent:      profile.UserAgent,
		AcceptLanguage: profile.AcceptLanguage,
		AcceptEncoding: profile.AcceptEncoding,
	}
}

// InvalidateBrowserIdentity drops the cached browser surface for one account.
func InvalidateBrowserIdentity(accountKey string) {
	infraegress.InvalidateStickyBrowserProfile(accountKey)
}

func accountBrowserKey(credentialID uint64, egressIdentity string) string {
	if value := strings.TrimSpace(egressIdentity); value != "" {
		return value
	}
	if credentialID == 0 {
		return ""
	}
	return fmt.Sprintf("console_%d", credentialID)
}
