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

	"github.com/pkg/errors"
	"goodrain.com/cloud-adaptor/internal/model"
	"goodrain.com/cloud-adaptor/pkg/bcode"
	"goodrain.com/cloud-adaptor/pkg/util/uuidutil"
	"gorm.io/gorm"
)

// RKEClusterRepo -
type RKEClusterRepo struct {
	DB *gorm.DB `inject:""`
}

// NewRKEClusterRepo creates a new RKEClusterRepository.
func NewRKEClusterRepo(db *gorm.DB) RKEClusterRepository {
	return &RKEClusterRepo{DB: db}
}

// Create create an event
func (t *RKEClusterRepo) Create(te *model.RKECluster) error {
	if te.Name == "" || te.EnterpriseID == "" {
		return fmt.Errorf("rke cluster name or eid can not be empty")
	}
	var old model.RKECluster
	if err := t.DB.Where("name=? and eid=?", te.Name, te.EnterpriseID).Take(&old).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// not found error, create new
			if te.ClusterID == "" {
				te.ClusterID = uuidutil.NewUUID()
			}
			if err := t.DB.Save(te).Error; err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return errors.WithStack(bcode.ErrRKEClusterExists)
}

// Update -
func (t *RKEClusterRepo) Update(te *model.RKECluster) error {
	return t.DB.Save(te).Error
}

// GetCluster -
func (t *RKEClusterRepo) GetCluster(eid, name string) (*model.RKECluster, error) {
	var rc model.RKECluster
	if err := t.DB.Where("eid=? and (name=? or clusterID=?)", eid, name, name).Take(&rc).Error; err != nil {
		return nil, err
	}
	return &rc, nil
}

// ListCluster -
func (t *RKEClusterRepo) ListCluster(eid string) ([]*model.RKECluster, error) {
	var list []*model.RKECluster
	if err := t.DB.Where("eid=?", eid).Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ListInitCluster - 查找到所有正在初始化的集群
func (t *RKEClusterRepo) ListInitCluster() []*model.RKECluster {
	var list = make([]*model.RKECluster, 0)
	if err := t.DB.Where("stats=?", "initial").Order("id desc").Find(&list).Error; err != nil {
		return list
	}
	return list
}

// DeleteCluster delete cluster
func (t *RKEClusterRepo) DeleteCluster(eid, name string) error {
	var rc model.RKECluster
	if err := t.DB.Where("eid=? and (name=? or clusterID=?)", eid, name, name).Delete(&rc).Error; err != nil {
		return err
	}
	return nil
}
