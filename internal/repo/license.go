// RAINBOND, Application Management Platform
// Copyright (C) 2020-2020 Goodrain Co., Ltd.

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

package repo

import (
	"fmt"
	"goodrain.com/cloud-adaptor/internal/model"
	"gorm.io/gorm"
)

// LicenseRepo license repo
type LicenseRepo struct {
	DB *gorm.DB
}

// NewLicenseRepo new license repo
func NewLicenseRepo(db *gorm.DB) LicenseRepository {
	return &LicenseRepo{DB: db}
}

func (l *LicenseRepo) GetFirstEnterprise() (*model.Enterprise, error) {
	var ent model.Enterprise
	if err := l.DB.Raw("select * from tenant_enterprise limit 1").Scan(&ent).Error; err != nil {
		return nil, err
	}
	return &ent, nil
}

func (l *LicenseRepo) GetRegionsByEID(eid string) ([]*model.Region, error) {
	var regions []*model.Region
	sql := fmt.Sprintf("select * from region_info where status = 1 and enterprise_id = '%s'", eid)
	if err := l.DB.Raw(sql).Scan(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}
