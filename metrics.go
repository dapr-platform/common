package common

type MetricPanel struct {
	Row   int    `json:"row"`
	Col   int    `json:"col"`
	Title string `json:"title"`
	Query string `json:"query"`
	Type  string `json:"type"` //trend,bar,current,stack,etc..
}

var DefaultmetricsPanel = []MetricPanel{
	{Row: 0, Col: 0, Title: "go协程数", Query: "go_goroutines{ident=\"${HOST}\",instance=\"${NAME}:80\"}", Type: "trend"},
	{Row: 0, Col: 1, Title: "go gc秒数", Query: "go_gc_duration_seconds{ident=\"${HOST}\",instance=\"${NAME}:80\"}", Type: "trend"},
	{Row: 0, Col: 2, Title: "go线程数", Query: "go_threads{ident=\"${HOST}\",instance=\"${NAME}:80\"}", Type: "trend"},
	{Row: 0, Col: 3, Title: "go分配对象数", Query: "go_memstats_heap_objects{ident=\"${HOST}\",instance=\"${NAME}:80\"}", Type: "trend"},
}
