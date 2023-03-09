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

package frontend

import (
	"crypto/rsa"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/balancereader"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/contacts"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/ledgerwriter"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/transactionhistory"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/userservice"
	"github.com/golang-jwt/jwt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	//go:embed static/*
	staticFS embed.FS

	validEnvs = []string{"local", "gcp"}
)

type platformDetails struct {
	css      string
	provider string
}

func (plat *platformDetails) setPlatformDetails(env string) {
	if env == "gcp" {
		plat.provider = "Google Cloud"
		plat.css = "gcp-platform"
	} else {
		plat.provider = "local"
		plat.css = "local"
	}
}

// ServerConfig contains configuration options for the server.
type ServerConfig struct {
	publicKey       *rsa.PublicKey
	localRoutingNum string
	bankName        string
	backendTimeout  time.Duration
	clusterName     string
	podName         string
	podZone         string
}

// Server is the application frontend.
type Server struct {
	handler  http.Handler
	root     weaver.Instance
	platform platformDetails
	hostname string
	config   ServerConfig

	balanceReader      balancereader.T
	contacts           contacts.T
	ledgerWriter       ledgerwriter.T
	transactionHistory transactionhistory.T
	userService        userservice.T
}

// NewServer returns a new application frontend.
func NewServer(root weaver.Instance, publicKeyPath, localRoutingNum, bankName string, backendTimeout time.Duration) (*Server, error) {
	// Setup the services.
	balanceReader, err := weaver.Get[balancereader.T](root)
	if err != nil {
		return nil, err
	}
	root.Logger().Debug("Initialized component: balancereader")
	contacts, err := weaver.Get[contacts.T](root)
	if err != nil {
		return nil, err
	}
	root.Logger().Debug("Initialized component: contacts")
	ledgerWriter, err := weaver.Get[ledgerwriter.T](root)
	if err != nil {
		return nil, err
	}
	root.Logger().Debug("Initialized component: ledgerwriter")
	transactionHistory, err := weaver.Get[transactionhistory.T](root)
	if err != nil {
		return nil, err
	}
	root.Logger().Debug("Initialized component: transactionhistory")
	userService, err := weaver.Get[userservice.T](root)
	if err != nil {
		return nil, err
	}
	root.Logger().Debug("Initialized component: userservice")

	// Find out where we're running.
	var env = os.Getenv("ENV_PLATFORM")
	// Only override from env variable if set + valid env.
	if env == "" || !stringinSlice(validEnvs, env) {
		root.Logger().Debug("ENV_PLATFORM is either empty or invalid")
		env = "local"
	}
	// Autodetect GCP.
	addrs, err := net.LookupHost("metadata.google.internal.")
	if err == nil && len(addrs) >= 0 {
		root.Logger().Debug("Detected Google metadata server, setting ENV_PLATFORM to GCP.", "address", addrs)
		env = "gcp"
	}
	root.Logger().Debug("ENV_PLATFORM", "platform", env)
	platform := platformDetails{}
	platform.setPlatformDetails(strings.ToLower(env))
	hostname, err := os.Hostname()
	if err != nil {
		root.Logger().Debug(`cannot get hostname for frontend: using "unknown"`)
		hostname = "unknown"
	}

	metadataServer := os.Getenv("METADATA_SERVER")
	if metadataServer == "" {
		metadataServer = "metadata.google.internal"
	}
	metadataURL := fmt.Sprintf("http://%s/computeMetadata/v1/", metadataServer)
	metadataHeaders := http.Header{}
	metadataHeaders.Set("Metadata-Flavor", "Google")
	clusterName := getClusterName(metadataURL, metadataHeaders)
	podName, err := os.Hostname()
	podZone := getPodZone(metadataURL, metadataHeaders)

	var c ServerConfig
	pubKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %v", err)
	}
	c.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %v", err)
	}
	c.localRoutingNum = localRoutingNum
	c.backendTimeout = backendTimeout
	c.bankName = bankName
	c.clusterName = clusterName
	c.podName = podName
	c.podZone = podZone

	// Create the server.
	s := &Server{
		root:               root,
		platform:           platform,
		hostname:           hostname,
		config:             c,
		balanceReader:      balanceReader,
		contacts:           contacts,
		ledgerWriter:       ledgerWriter,
		transactionHistory: transactionHistory,
		userService:        userService,
	}

	// Setup the handler.
	staticHTML, err := fs.Sub(fs.FS(staticFS), "static")
	if err != nil {
		return nil, err
	}
	r := http.NewServeMux()

	// Helper that adds a handler with HTTP metric instrumentation.
	instrument := func(label string, fn func(http.ResponseWriter, *http.Request), methods []string) http.Handler {
		allowed := map[string]struct{}{}
		for _, method := range methods {
			allowed[method] = struct{}{}
		}
		handler := func(w http.ResponseWriter, r *http.Request) {
			if _, ok := allowed[r.Method]; len(allowed) > 0 && !ok {
				msg := fmt.Sprintf("method %q not allowed", r.Method)
				http.Error(w, msg, http.StatusMethodNotAllowed)
				return
			}
			fn(w, r)
		}
		return weaver.InstrumentHandlerFunc(label, handler)
	}

	const get = http.MethodGet
	const post = http.MethodPost
	const head = http.MethodHead
	r.Handle("/", instrument("root", s.rootHandler, []string{get, head}))
	r.Handle("/home/", instrument("home", s.homeHandler, []string{get, head}))
	r.Handle("/payment", instrument("payment", s.paymentHandler, []string{post}))
	r.Handle("/deposit", instrument("deposit", s.depositHandler, []string{post}))
	r.Handle("/login", instrument("login", s.loginHandler, []string{get, post}))
	r.Handle("/consent", instrument("consent", s.consentHandler, []string{get, post}))
	r.Handle("/signup", instrument("signup", s.signupHandler, []string{get, post}))
	r.Handle("/logout", instrument("logout", s.logoutHandler, []string{post}))
	r.Handle("/static/", weaver.InstrumentHandler("static", http.StripPrefix("/static", http.FileServer(http.FS(staticHTML)))))

	// No instrumentation of /healthz
	r.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	// Set handler and return.
	var handler http.Handler = r
	handler = newLogHandler(root, handler)         // add logging
	handler = otelhttp.NewHandler(handler, "http") // add tracing
	s.handler = handler

	return s, nil
}

func getClusterName(metadataURL string, metadataHeaders http.Header) string {
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		clusterName = "unknown"
	}
	req, err := http.NewRequest("GET", metadataURL+"instance/attributes/cluster-name", nil)
	if err != nil {
		return clusterName
	}
	req.Header = metadataHeaders

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return clusterName
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return clusterName
	} else {
		clusterNameBytes := make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(clusterNameBytes)
		if err != nil {
			return clusterName
		}
		clusterName = string(clusterNameBytes)
	}
	return clusterName
}

func getPodZone(metadataURL string, metadataHeaders http.Header) string {
	podZone := os.Getenv("POD_ZONE")
	if podZone == "" {
		podZone = "unknown"
	}
	req, err := http.NewRequest("GET", metadataURL+"instance/zone", nil)
	if err != nil {
		return podZone
	}
	req.Header = metadataHeaders

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return podZone
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		podZoneBytes := make([]byte, resp.ContentLength)
		_, err = resp.Body.Read(podZoneBytes)
		if err != nil {
			return podZone
		}
		podZoneSplit := strings.Split(string(podZoneBytes), "/")
		if len(podZoneSplit) >= 4 {
			podZone = podZoneSplit[3]
		}
	}
	return podZone
}

func (s *Server) Run(localAddr string) error {
	lis, err := s.root.Listener("bank", weaver.ListenerOptions{LocalAddress: localAddr})
	if err != nil {
		return err
	}
	s.root.Logger().Debug("Frontend available", "addr", lis)
	return http.Serve(lis, s.handler)
}
