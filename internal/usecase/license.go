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

package usecase

import (
	"github.com/sirupsen/logrus"
	"goodrain.com/cloud-adaptor/internal/region"
	"goodrain.com/cloud-adaptor/internal/repo"
	licenseutil "goodrain.com/cloud-adaptor/pkg/util/license"
	"goodrain.com/cloud-adaptor/pkg/util/timeutil"
	"time"
)

type licenseCache struct {
	license licenseutil.LicenseResp
	timeOut time.Time
}

func (c *licenseCache) IsExpired() bool {
	return c.timeOut.Before(time.Now())
}

// LicenseUsecase license usecase
type LicenseUsecase struct {
	LicenseRepo    repo.LicenseRepository
	regionLicenses map[string]licenseCache
}

// NewLicenseUsecase new license usecase
func NewLicenseUsecase(licenseRepo repo.LicenseRepository) *LicenseUsecase {
	return &LicenseUsecase{
		LicenseRepo:    licenseRepo,
		regionLicenses: make(map[string]licenseCache),
	}
}

// GetLicense -
func (l *LicenseUsecase) GetLicense() *licenseutil.AllLicense {
	consoleLicense := licenseutil.ReadLicense()
	allLicense := &licenseutil.AllLicense{}
	// handle console license
	if consoleLicense == nil {
		return allLicense
	}
	allLicense.HaveLicense = true
	allLicense.EndTime = consoleLicense.EndTime
	allLicense.RegionNums = consoleLicense.Cluster
	if consoleLicense.EndTime == "" {
		allLicense.IsPermanent = true
	}
	if !allLicense.IsPermanent && timeutil.JudgeTimeIsExpired(consoleLicense.EndTime) {
		allLicense.IsExpired = true
	}
	// get region licenses
	regionLicenses, err := l.GetRegionLicenses()
	if err != nil {
		return allLicense
	}
	allLicense.RegionLicenses = regionLicenses
	return allLicense
}

// GetRegionLicenses -
func (l *LicenseUsecase) GetRegionLicenses() ([]*licenseutil.LicenseResp, error) {
	ent, err := l.LicenseRepo.GetFirstEnterprise()
	if err != nil {
		return nil, err
	}
	regions, err := l.LicenseRepo.GetRegionsByEID(ent.EnterpriseID)
	if err != nil {
		return nil, err
	}
	var licenses []*licenseutil.LicenseResp
	for _, rg := range regions {
		if cache, ok := l.regionLicenses[rg.RegionName]; ok {
			if cache.IsExpired() {
				delete(l.regionLicenses, rg.RegionName)
			} else {
				logrus.Infof("get regin [%s] license by cache", rg.RegionName)
				licenses = append(licenses, &cache.license)
				continue
			}
		}
		rainbondClient, err := region.NewRegion(region.APIConf{
			Endpoints: []string{rg.URL},
			Token:     rg.Token,
			Cacert:    rg.SSlCaCert,
			Cert:      rg.CertFile,
			CertKey:   rg.KeyFile,
		})
		if err != nil {
			logrus.Errorf("new rainbond client failed %v", err)
			continue
		}
		license, err := rainbondClient.License().Get()
		if err != nil && license == nil {
			logrus.Errorf("get rainbond license failed %v", err)
			continue
		}
		license.RegionName = rg.RegionName
		licenses = append(licenses, license)
		l.regionLicenses[rg.RegionName] = licenseCache{license: *license, timeOut: time.Now().Add(time.Hour * 24)}
	}
	return licenses, nil
}
