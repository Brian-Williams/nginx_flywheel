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
		Kvs: []*mvccpb.KeyValue{&kv},
	}

	directives := kvToDirectives(request)

	key := directives[0].Directive
	if key != "Cat" {
		t.Errorf("Incorrect new key, expected 'Cat' got: %s", key)
	}
	value := directives[0].Args[0]
	if value != "Dog" {
		t.Errorf("Incorrect new value, expected 'Dog' got: %s", value)
	}
}
