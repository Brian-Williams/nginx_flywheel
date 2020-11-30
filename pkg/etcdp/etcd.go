package etcdp

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/Brian-Williams/nginx_flywheel/pkg"

	"github.com/coreos/etcd/clientv3"
)

// Etcd3Provider is a OverrideProvider for etcd
type Etcd3Provider struct {
	_ struct{}
	*clientv3.Client
	// LStrip is the prefix strip for NGINX config location
	LStrip string
}

var _ flywheel.OverrideProvider = (*Etcd3Provider)(nil)

// New creates an Etcd3Provider
func New(config clientv3.Config, lstrip string) (*Etcd3Provider, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &Etcd3Provider{Client: cli, LStrip: lstrip}, nil
}

// Override satisfies the OverrideProvider interface
func (e *Etcd3Provider) Override(directive, path string) ([]string, error) {
	key := e.DirectiveKey(directive, path)

	// Context should be configured in New, so that ctx doesn't infect every method
	r, err := e.Client.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}
	values := make([]string, len(r.Kvs))
	for i, v := range r.Kvs {
		values[i] = string(v.Value)
	}
	return values, nil
}

// DirectiveKey produces a key from a directive and NGINX filepath
//
// For example directive listen with path /etc/nginx/nginx.conf and LStrip /etc/nginx would produce /nginx/listen
func (e *Etcd3Provider) DirectiveKey(directive, path string) string {
	return strings.TrimPrefix(filepath.Clean(strings.TrimSuffix(path, filepath.Ext(path))), e.LStrip) + "/" + directive
}
