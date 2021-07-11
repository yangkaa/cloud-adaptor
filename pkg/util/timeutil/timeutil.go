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

package timeutil

import "time"

// JudgeTimeIsExpired -
func JudgeTimeIsExpired(endTime string) bool {
	loc, _ := time.LoadLocation("Local")
	theEndTime, err := time.ParseInLocation("2006-01-02 15:04:05", endTime, loc)
	if err != nil {
		return false
	}
	if theEndTime.Before(time.Now()){
		return true
	}
	return false
}
