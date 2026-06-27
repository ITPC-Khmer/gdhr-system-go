// Route absolute http(s) image URLs through the backend image proxy so the
// browser never contacts the external host directly. Relative/empty values and
// data URLs are returned unchanged.
const base = import.meta.env.VITE_API_BASE_URL || '/api'

export function proxied(url) {
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) {
    return `${base}/image-proxy?url=${encodeURIComponent(url)}`
  }
  return url
}
