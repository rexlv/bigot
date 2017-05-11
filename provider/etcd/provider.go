package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"errors"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

// Provider is a wrapper around the etcd client
type Provider struct {
	client.KeysAPI
	keys []string
}

// NewProvider returns an *etcd.Client with a connection to named machines.
func NewProvider(machines []string, cert, key, caCert string, basicAuth bool, username string, password string, keys ...string) (p *Provider, err error) {
	var c client.Client
	var kapi client.KeysAPI
	var transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	cfg := client.Config{
		Endpoints:               machines,
		HeaderTimeoutPerRequest: time.Duration(3) * time.Second,
	}

	if basicAuth {
		cfg.Username = username
		cfg.Password = password
	}

	if caCert != "" {
		var certBytes []byte
		certBytes, err = ioutil.ReadFile(caCert)
		if err != nil {
			return
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
	}

	if cert != "" && key != "" {
		var tlsCert tls.Certificate
		tlsCert, err = tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
	}

	transport.TLSClientConfig = tlsConfig
	cfg.Transport = transport

	c, err = client.New(cfg)
	if err != nil {
		return
	}

	kapi = client.NewKeysAPI(c)
	return &Provider{KeysAPI: kapi, keys: keys}, nil
}

func (c *Provider) Read() (config map[string]interface{}, err error) {
	config = make(map[string]interface{})
	for _, key := range c.keys {
		if err = c.readKey(key, config); err != nil {
			break
		}
	}
	return
}

func (c *Provider) readKey(key string, config map[string]interface{}) (err error) {
	var resp *client.Response
	resp, err = c.Get(context.Background(), key, &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	})

	if err != nil {
		return
	}

	if resp.Node == nil {
		return errors.New("no found node")
	}

	c.walkNodes(resp.Node, config)
	return
}

func (c *Provider) walkNodes(node *client.Node, config map[string]interface{}) {
	if !node.Dir {
		fmt.Println(node.Key, " : ", node.Value)
		config[node.Key] = node.Value
		return
	}

	for _, n := range node.Nodes {
		c.walkNodes(n, config)
	}
}
