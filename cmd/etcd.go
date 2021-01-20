/*
Copyright © 2020 Brian Williams

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/Brian-Williams/nginx_flywheel/pkg/etcdp"

	"github.com/Brian-Williams/nginx_flywheel/pkg"
	"github.com/aluttik/go-crossplane"
	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	endpoints []string
	lstrip    string

	// etcdCmd represents the etcd command
	etcdCmd = &cobra.Command{
		Use:   "etcd",
		Short: "Rewrite an NGINX file using etcd keys as a variable provider",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().
				Str("path", sourcePath).
				Msg("Parsing file")
			payload, err := crossplane.Parse(sourcePath, &crossplane.ParseOptions{ParseComments: true})
			if err != nil {
				msg := "failed to parse source file"
				log.Err(err).Str("file", sourcePath).Msg(msg)
				return fmt.Errorf(msg+": %w", err)
			}

			log.Print("Replacing directive keys from etcd")
			client, err := clientv3.New(clientv3.Config{Endpoints: endpoints})
			if err != nil {
				msg := "invalid NGINX configuration"
				log.Err(err).Strs("endpoints", endpoints).Msg(msg)
				return fmt.Errorf(msg)
			}
			overrider := &etcdp.Etcd3Provider{Client: client, LStrip: lstrip}

			err = flywheel.OverridePayload(context.Background(), payload, overrider)
			if err != nil {
				msg := "overriding NGINX JSON failed"
				log.Err(err).Msg(msg)
				return fmt.Errorf(msg)
			}

			log.Print("Writing payload")
			err = flywheel.WritePayload(payload, &crossplane.BuildOptions{})
			if err != nil {
				msg := "failed to write to output"
				log.Err(err).Msg(msg)
				return fmt.Errorf(msg)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(etcdCmd)
	// etcd flags
	etcdCmd.PersistentFlags().StringSliceVar(&endpoints, "endpoint", nil, "etcd endpoints")
	etcdCmd.PersistentFlags().StringVar(&lstrip, "lstrip", "/etc", "prefix to strip from absolute path to produce etcd key")
}
