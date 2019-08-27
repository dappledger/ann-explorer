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
