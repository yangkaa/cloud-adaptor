// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

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

package util

import (
	"math/rand"
	"net"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString create rand string
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// GetIPByURL get ip by url
func GetIPByURL(u string) string {
	url, _ := url.Parse(u)
	if url != nil {
		hostIP := strings.Split(url.Host, ":")
		var ip string
		if len(hostIP) > 0 {
			ip = hostIP[0]
		}
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}

// CheckCommandExists checks if a command exists in the system
func CheckCommandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// CheckPortInUse checks if a port is in use on 127.0.0.1
func CheckPortInUse(port int) bool {
	address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	conn, err := net.Listen("tcp", address)
	if err != nil {
		return true // Port is in use
	}
	conn.Close()
	return false // Port is not in use
}
