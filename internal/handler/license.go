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

package handler

import (
	"github.com/gin-gonic/gin"
	"goodrain.com/cloud-adaptor/internal/usecase"
	"goodrain.com/cloud-adaptor/pkg/util/ginutil"
)

// LicenseHandler -
type LicenseHandler struct {
	LicenseUsecase usecase.LicenseUsecase
}

// NewLicenseHandler new license handler
func NewLicenseHandler(licenseUcase usecase.LicenseUsecase) *LicenseHandler {
	return &LicenseHandler{LicenseUsecase: licenseUcase}
}

// GetLicense -
func (l *LicenseHandler) GetLicense(ctx *gin.Context) {
	allLicense := l.LicenseUsecase.GetLicense()
	ginutil.JSON(ctx, allLicense)
}