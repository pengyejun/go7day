package remote

import (
	"bytes"
	"fmt"
	"gee-cache/geecache"
	"gee-cache/geecachepb"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"net/url"
	"sync"
)

type httpGetter struct {
	baseURL string
}

var bufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func (h *httpGetter) Get(in *geecachepb.Request, out *geecachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)

	res, err := http.Get(u)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("remote returned: %v", res.Status)
	}
	resBytes := bufferPool.Get().(*bytes.Buffer)
	resBytes.Reset()
	defer bufferPool.Put(resBytes)
	_, err = io.Copy(resBytes, res.Body)
	if err != nil {
		return fmt.Errorf("reading responce body: %v", err)
	}
	if err = proto.Unmarshal(resBytes.Bytes(), out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ geecache.PeerGetter = (*httpGetter)(nil)