package dns

import "sync"

// There are 3 buckets of the compression map pools. Each each pool needs to be roughly cache
// elements of the same size. Here we define a small, medium and large pool. When we are packing a
// message we don't know the resulting size (yet), so we estimate the size based on the number of records
// we are packing

var (

	// compressionPackSmallPool caches maps for messages that have <= 50 RRs in them.
	compressionPackSmallPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]uint16)
		},
	}
	// map for the Len function.
	compressionLenSmallPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]struct{})
		},
	}
	// compressionPackMediumPool caches maps for messages that have > 50 and <= 250 RRs in them.
	compressionPackMediumPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]uint16)
		},
	}
	// map for the Len function.
	compressionLenMediumPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]struct{})
		},
	}
	// compressionPackLargePool caches maps for messages that have > 250 RRs in them.
	compressionPackLargePool = sync.Pool{
		New: func() interface{} {
			return make(map[string]uint16)
		},
	}
	// map for the Len function.
	compressionLenLargePool = sync.Pool{
		New: func() interface{} {
			return make(map[string]struct{})
		},
	}
)

// compressionPackPool returns the correct pool based on the number of RRs in dns.
func compressionPackPool(dns *Msg) sync.Pool {
	rrs := len(dns.Question) + len(dns.Answer) + len(dns.Ns) + len(dns.Extra)
	switch {
	case rrs <= 50:
		return compressionPackSmallPool
	case rrs > 50 && rrs <= 250:
		return compressionPackMediumPool
	case rrs > 250:
		return compressionPackLargePool
	}
	// not reached
	return compressionPackSmallPool
}

// compressionLenPool returns the correct pool based on the number of RRs in dns.
func compressionLenPool(dns *Msg) sync.Pool {
	rrs := len(dns.Question) + len(dns.Answer) + len(dns.Ns) + len(dns.Extra)
	switch {
	case rrs <= 50:
		return compressionLenSmallPool
	case rrs > 50 && rrs <= 250:
		return compressionLenMediumPool
	case rrs > 250:
		return compressionLenLargePool
	}
	// not reached
	return compressionLenSmallPool
}
