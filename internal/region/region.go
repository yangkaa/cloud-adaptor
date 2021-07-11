// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package region

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	utilhttp "goodrain.com/cloud-adaptor/pkg/util/httputil"
)

var regionAPI, token string

var region Region

//AllTenant AllTenant
var AllTenant string

//Region region api
type Region interface {
	DoRequest(path, method string, body io.Reader, decode *utilhttp.ResponseBody) (int, error)
	License() LicenseInterface
}

//APIConf region api config
type APIConf struct {
	Endpoints []string `yaml:"endpoints"`
	Token     string   `yaml:"token"`
	AuthType  string   `yaml:"auth_type"`
	Cacert    string   `yaml:"client-ca-file"`
	Cert      string   `yaml:"tls-cert-file"`
	CertKey   string   `yaml:"tls-private-key-file"`
}

//NewRegion NewRegion
func NewRegion(c APIConf) (Region, error) {
	if region == nil {
		re := &regionImpl{
			APIConf: c,
		}
		if c.Cacert != "" && c.Cert != "" && c.CertKey != "" {
			pool := x509.NewCertPool()
			caCrt, err := ioutil.ReadFile(c.Cacert)
			if err != nil {
				logrus.Errorf("read ca file err: %s", err)
				return nil, err
			}
			pool.AppendCertsFromPEM(caCrt)
			cliCrt, err := tls.LoadX509KeyPair(c.Cert, c.CertKey)
			if err != nil {
				logrus.Errorf("Loadx509keypair err: %s", err)
				return nil, err
			}
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      pool,
					Certificates: []tls.Certificate{cliCrt},
				},
			}
			re.Client = &http.Client{
				Transport: tr,
				Timeout:   15 * time.Second,
			}
		} else {
			re.Client = http.DefaultClient
		}
		region = re
	}
	return region, nil
}

//GetRegion GetRegion
func GetRegion() Region {
	return region
}

type regionImpl struct {
	APIConf
	Client *http.Client
}

func (r *regionImpl) GetEndpoint() string {
	return r.Endpoints[0]
}

//DoRequest do request
func (r *regionImpl) DoRequest(path, method string, body io.Reader, decode *utilhttp.ResponseBody) (int, error) {
	request, err := http.NewRequest(method, r.GetEndpoint()+path, body)
	if err != nil {
		return 500, err
	}
	request.Header.Set("Content-Type", "application/json")
	if r.Token != "" {
		request.Header.Set("Authorization", "Token "+r.Token)
	}
	res, err := r.Client.Do(request)
	if err != nil {
		return 500, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if decode != nil {
		if err := json.NewDecoder(res.Body).Decode(decode); err != nil {
			return res.StatusCode, err
		}
	}
	return res.StatusCode, err
}
