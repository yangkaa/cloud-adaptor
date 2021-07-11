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

package http

import (
	"net/url"
	"sync"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	var once sync.Once
	once.Do(func() {
		validate = validator.New()
	})
}

// ErrBadRequest -
type ErrBadRequest struct {
	err error
}

func (e ErrBadRequest) Error() string {
	return e.err.Error()
}

// Result represents a response for restful api.
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//ResponseBody api返回数据格式
type ResponseBody struct {
	ValidationError url.Values  `json:"validation_error,omitempty"`
	Msg             string      `json:"msg,omitempty"`
	Bean            interface{} `json:"bean,omitempty"`
	List            interface{} `json:"list,omitempty"`
	//数据集总数
	ListAllNumber int `json:"number,omitempty"`
	//当前页码数
	Page int `json:"page,omitempty"`
}
