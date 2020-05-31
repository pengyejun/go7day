package geecache

import "gee-cache/geecachepb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *geecachepb.Request, out *geecachepb.Response) error
}


