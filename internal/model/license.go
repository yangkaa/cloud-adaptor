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

package model

// Enterprise -
type Enterprise struct {
	EnterpriseID string `gorm:"column:enterprise_id"`
}

type Region struct {
	EnterpriseID string `gorm:"column:enterprise_id"`
	RegionName   string `gorm:"column:region_name"`
	SSlCaCert    string `gorm:"column:ssl_ca_cert"`
	CertFile     string `gorm:"column:cert_file"`
	KeyFile      string `gorm:"column:key_file"`
	URL          string `gorm:"column:url"`
	Token        string `gorm:"column:token"`
}
