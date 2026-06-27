package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// maxImageBytes caps a single proxied image to guard against memory abuse.
const maxImageBytes = 10 << 20 // 10 MB

var (
	imageProxyHosts  = map[string]bool{}
	imageProxyClient *http.Client
)

// InitImageProxy configures the allowlist of hosts the image proxy may fetch
// from, plus a redirect-validating HTTP client. Called once at startup.
func InitImageProxy(hosts []string) {
	imageProxyHosts = map[string]bool{}
	for _, h := range hosts {
		if h = strings.TrimSpace(strings.ToLower(h)); h != "" {
			imageProxyHosts[h] = true
		}
	}
	imageProxyClient = &http.Client{
		Timeout: 20 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Re-validate every redirect hop so a redirect can't escape the
			// allowlist (SSRF protection). Stop following on violation.
			if len(via) >= 5 || !hostAllowed(req.URL.Hostname()) {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
}

func hostAllowed(host string) bool {
	return imageProxyHosts[strings.ToLower(host)]
}

// ImageProxy fetches an allowlisted remote image and streams it back, so the
// browser never contacts the external host directly. The host allowlist makes
// this safe to expose unauthenticated (it can only ever return those images),
// which lets it be used straight from an <img src> tag.
//
//	GET /api/image-proxy?url=<absolute http(s) image URL>
func ImageProxy(c *gin.Context) {
	if imageProxyClient == nil {
		InitImageProxy([]string{"police-gdhr.interior.gov.kh"})
	}

	raw := c.Query("url")
	if raw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing url"})
		return
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid url"})
		return
	}
	if !hostAllowed(u.Hostname()) {
		c.JSON(http.StatusForbidden, gin.H{"message": "image host not allowed"})
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "failed to build request"})
		return
	}

	resp, err := imageProxyClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Status(http.StatusBadGateway)
		return
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(strings.ToLower(ct), "image/") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"message": "remote resource is not an image"})
		return
	}

	c.Header("Content-Type", ct)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, io.LimitReader(resp.Body, maxImageBytes))
}
