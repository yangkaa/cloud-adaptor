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
	"fmt"
	"github.com/sirupsen/logrus"
	"goodrain.com/cloud-adaptor/internal/region"
	"goodrain.com/cloud-adaptor/internal/repo"
	licenseutil "goodrain.com/cloud-adaptor/pkg/util/license"
	"goodrain.com/cloud-adaptor/pkg/util/timeutil"
	"os"
	"time"
)

// LicenseUsecase license usecase
type LicenseUsecase struct {
	LicenseRepo repo.LicenseRepository
}

// NewLicenseUsecase new license usecase
func NewLicenseUsecase(licenseRepo repo.LicenseRepository) *LicenseUsecase {
	return &LicenseUsecase{
		LicenseRepo: licenseRepo,
	}
}

// GetLicense -
func (l *LicenseUsecase) GetLicense() *licenseutil.AllLicense {
	consoleLicense := licenseutil.ReadLicense()
	allLicense := &licenseutil.AllLicense{}
	// handle console license
	if consoleLicense == nil {
		allLicense.EndTime = time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
		return allLicense
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", consoleLicense.EndTime)
	if err != nil {
		logrus.Infof("parse end time failed %v", err)
		endTime = time.Now().Add(time.Hour * 24)
		consoleLicense.EndTime = endTime.Format("2006-01-02 15:04:05")
	}
	if endTime.IsZero() {
		allLicense.IsPermanent = true
	}
	if !allLicense.IsPermanent && timeutil.JudgeTimeIsExpired(consoleLicense.EndTime) {
		allLicense.EndTime = consoleLicense.EndTime
		allLicense.IsExpired = true
	}
	allLicense.RegionNums = consoleLicense.Cluster
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
		if err := l.generateCert(rg.RegionName, rg.SSlCaCert, rg.CertFile, rg.KeyFile); err != nil {
			continue
		}
		rainbondClient, err := region.NewRegion(region.APIConf{
			Endpoints: []string{rg.URL},
			Token:     rg.Token,
			Cacert:    fmt.Sprintf("./certs/%s_ca.pem", rg.RegionName),
			Cert:      fmt.Sprintf("./certs/%s_client.pem", rg.RegionName),
			CertKey:   fmt.Sprintf("./certs/%s_client_key.pem", rg.RegionName),
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
	}
	return licenses, nil
}

func (l *LicenseUsecase) generateCert(regionName, caCert, cert, certKey string) error {
	err := os.MkdirAll("./certs", os.ModePerm)
	if err != nil {
		logrus.Errorf("create cert directly failed %v", err)
		return err
	}
	fileNames := map[string]string{
		fmt.Sprintf("./certs/%s_ca.pem", regionName):         caCert,
		fmt.Sprintf("./certs/%s_client.pem", regionName):     cert,
		fmt.Sprintf("./certs/%s_client_key.pem", regionName): certKey,
	}
	for fileName, content := range fileNames {
		f, err := os.Create(fileName)
		if err != nil {
			logrus.Errorf("create ca cert failed %v", err)
			return err
		}
		_, err = f.WriteString(content)
		if err != nil {
			logrus.Errorf("write cert %s failed %v", fileName, err)
			f.Close()
			return err
		}
		err = f.Close()
		if err != nil {
			logrus.Errorf("close cert file %s failed %v", fileName, err)
			return err
		}
	}
	return nil
}
