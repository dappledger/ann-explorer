// Copyright © 2017 ZhongAn Technology
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import "strconv"

const pageSize = 50

type Pagination struct {
	Total   int `json:"total"`
	Current int `json:"current"`
}

type PaginationResult struct {
	Data interface{} `json:"data"`
	Page Pagination  `json:"page"`
}

// -1 => latest
func getPage(wc *WebController) int {
	p := wc.GetString("p")
	if p == "latest" || p == "" {
		return 1
	}
	i, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return 1
	}
	return int(i)
}

func totalPage(count int) int {

	c := count / DisplayNum
	if count%DisplayNum != 0 {
		c += 1
	}
	return c
}

func calcFromTo(p, count int) (from, to int) {

	from = count - p*DisplayNum
	if p == 1 {
		to = count
		if from < 0 {
			from = 0
		}
		return
	}
	to = from + DisplayNum
	return
}
