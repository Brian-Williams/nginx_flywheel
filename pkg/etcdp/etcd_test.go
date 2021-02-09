package etcdp

import (
	"testing"

	"github.com/coreos/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
)

var (
	getResponse = clientv3.GetResponse{}
)

func TestKVToDirectives(t *testing.T) {
	kv := mvccpb.KeyValue{
		Key:   []byte("NEWCat"),
		Value: []byte("Dog"),
	}
	request := &clientv3.GetResponse{
		Kvs: []*mvccpb.KeyValue{&kv, &kv},
	}

	directives := kvToDirectives(request)

	for _, d := range directives {
		if d.Directive != "Cat" {
			t.Errorf("Incorrect new directive, expected 'Cat' got: %s", d.Directive)
		}
		if d.Args[0] != "Dog" {
			t.Errorf("Incorrect new args, expected 'Dog' to be first item in args got: %s", d.Args)
		}
	}
}
