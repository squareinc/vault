package transit

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/helper/keysutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func (b *backend) pathCacheConfig() *framework.Path {
	return &framework.Path{
		Pattern: "cache-config",
		Fields: map[string]*framework.FieldSchema{
			"cache-type": &framework.FieldSchema{
				Type:     framework.TypeString,
				Required: true,
				Description: `
Type of cache to use. Currently "syncmap" and "lru" are supported.
`,
			},

			"cache-size": &framework.FieldSchema{
				Type: framework.TypeInt,
				Description: `
Size of cache for a cache type that accepts a size. This is required for cache types
that accept a size and currently applies only to "lru" cache type.
`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.CreateOperation: b.pathCacheConfigWrite,
			logical.UpdateOperation: b.pathCacheConfigWrite,
			logical.ReadOperation:   b.pathCacheConfigRead,
		},

		HelpSynopsis:    pathCacheConfigHelpSyn,
		HelpDescription: pathCacheConfigHelpDesc,
	}
}

func (b *backend) pathCacheConfigWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	// get target cacheType
	cacheTypeStr := d.Get("cache-type").(string)
	cacheSize := d.Get("cache-size").(int)
	var cacheType keysutil.CacheType
	switch cacheTypeStr {
	case "syncmap":
		cacheType = keysutil.SyncMap
	case "lru":
		cacheType = keysutil.LRU
	default:
		cacheType = keysutil.NotImplemented
	}

	// err if the requested cacheType has not been implemented
	if cacheType == keysutil.NotImplemented {
		return logical.ErrorResponse(fmt.Sprintf("unknown cache-type %s", cacheTypeStr)), logical.ErrInvalidRequest
	}

	// err if cacheType is lru but no cache-size was specified
	if cacheType == keysutil.LRU && cacheSize <= 0 {
		return logical.ErrorResponse("for lru cache-type, cache-size must be specified and be greater than zero"), logical.ErrInvalidRequest
	}

	// change the cache type
	if cacheType == keysutil.SyncMap {
		b.lm.ConvertCacheToSyncmap()
	}

	if cacheType == keysutil.LRU {
		b.lm.ConvertCacheToLRU(cacheSize)
	}

	return nil, nil
}

func (b *backend) pathCacheConfigRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var cacheType string
	switch b.lm.GetCacheType() {
	case keysutil.SyncMap:
		cacheType = "syncMap"
	case keysutil.LRU:
		cacheType = "lru"
	default:
		cacheType = "unknown"
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"cache-type": cacheType,
			"cache-size": b.lm.GetCacheSize(),
		},
	}

	return resp, nil
}

const pathCacheConfigHelpSyn = `Configure caching strategy`

const pathCacheConfigHelpDesc = `
This path is used to configure and query the caching strategy for the transit mount.
For cache-types that do not have an upperbound like "syncmap", 0 is returned.
`
