// RAINBOND, Application Management Platform
// Copyright (C) 2021-2021 Goodrain Co., Ltd.

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
	"fmt"
	"github.com/sirupsen/logrus"
	"goodrain.com/cloud-adaptor/pkg/util"
	utilhttp "goodrain.com/cloud-adaptor/pkg/util/httputil"
	licenseutil "goodrain.com/cloud-adaptor/pkg/util/license"
	"time"
	"context"
)

func (r *regionImpl) License() LicenseInterface {
	return &license{regionImpl: *r, prefix: "/license"}
}

type license struct {
	regionImpl
	prefix string
}

type LicenseInterface interface {
	Get(ctx context.Context) (*licenseutil.LicenseResp, *util.APIHandleError)
}

// Get -
func (l *license) Get(ctx context.Context) (*licenseutil.LicenseResp, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var lic licenseutil.LicenseResp
	res.Bean = &lic

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	code, err := l.DoRequest(ctx, l.prefix, "GET", nil, &res)
	if err != nil {
		logrus.Debugf("------> request failed %v", err)
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		logrus.Debugf("------> request code failed %v", code)
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get license code %d", code))
	}
	return &lic, nil
}
