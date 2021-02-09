package etcdp

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aluttik/go-crossplane"

	"github.com/Brian-Williams/nginx_flywheel/pkg"

	"github.com/coreos/etcd/clientv3"
)

// New is the directive prefix for adding a directive
var New = "NEW"

// Etcd3Provider is a OverrideProvider for etcd
type Etcd3Provider struct {
	_ struct{}
	*clientv3.Client
	// LStrip is the prefix strip for NGINX config location
	LStrip string
}

var _ flywheel.OverrideProvider = (*Etcd3Provider)(nil)

// Override satisfies the OverrideProvider interface
func (e *Etcd3Provider) Override(ctx context.Context, directive, path string) ([]string, error) {
	key := e.DirectiveKey(directive, path)

	r, err := e.Client.Get(ctx, key)
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

// NewDirectives gets the directives to be added for a given path
//
// The path is used as a new prefix search so must end with '/'
func (e *Etcd3Provider) NewDirectives(ctx context.Context, path string) ([]crossplane.Directive, error) {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	r, err := e.Get(ctx, path+New, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get new directives: %w", err)
	}

	directives := kvToDirectives(r)

	return directives, nil
}

func kvToDirectives(r *clientv3.GetResponse) []crossplane.Directive {
	directives := make([]crossplane.Directive, len(r.Kvs))
	for i, v := range r.Kvs {
		d := strings.TrimPrefix(string(v.Key), New)
		directive := crossplane.Directive{
			Directive: d,
			Args:      []string{string(v.Value)},
		}
		directives[i] = directive
	}

	return directives
}
