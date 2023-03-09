// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main implements a demo banking application called Bank of Anthos.
//
// This application is a forked version of Google Cloud's Bank of Anthos
// app [1], with the following changes:
//   - It is written entirely in Go.
//   - It is written as a single Service Weaver application.
//   - It is written to use Service Weaver specific logging/tracing/monitoring.
//
// [1]: https://github.com/GoogleCloudPlatform/bank-of-anthos
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/frontend"
)

//go:generate weaver generate ./...

var (
	localAddr       = flag.String("local_addr", ":29471", "Local address")
	publicKeyPath   = flag.String("public_key_path", "/tmp/.ssh/jwtRS256.key.pub", "Path to the public key used for decrypting JWTs")
	localRoutingNum = flag.String("local_routing_num", "883745000", "The local routing number")
	backendTimeout  = flag.Duration("backend_timeout", 4*time.Second, "Timeout of calls to backend services")
	bankName        = flag.String("bank_name", "Bank of Anthos", "Name of the bank")
)

func main() {
	flag.Parse()
	root := weaver.Init(context.Background())
	server, err := frontend.NewServer(root, *publicKeyPath, *localRoutingNum, *bankName, *backendTimeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating frontend: ", err)
		os.Exit(1)
	}
	if err := server.Run(*localAddr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
