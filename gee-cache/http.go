package gee_cache


const defaultBasePath = "/_geecache/"

type HTTPPool struct {
	this string
	basePath string
}

func NewHTTPPool(this string) *HTTPPool {
	return &HTTPPool{
		this:     this,
		basePath: defaultBasePath,
	}
}