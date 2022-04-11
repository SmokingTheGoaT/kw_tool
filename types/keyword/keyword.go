package keyword

import (
	"kw_tool/util/enums"
	"strings"
)

//Map holds a keyword object for caching.
type Map struct {
	m  map[string]int
	ct *countTracer
	c  int
	p  enums.Platform
}

//Init is very important to initialize the keyword map before inserting
func (mp *Map) Init(p enums.Platform) {
	mp.m = make(map[string]int)
	mp.ct = &countTracer{}
	mp.c = 0
	mp.p = p
}

//Insert will allow for a key insert into the map if not already one
func (mp *Map) Insert(key string) {
	if _, found := mp.m[key]; !found {
		mp.parse(key)
		mp.m[key] = mp.c
		mp.c++
	}
}

//Platform will return a string of type platform
func (mp *Map) Platform() enums.Platform {
	return mp.p
}

//Len will return the map length
func (mp *Map) Len() int {
	return mp.c
}

//parse will parse  a key into an array of strings and keep count of frequency of words
func (mp *Map) parse(key string) {
	arr := strings.Fields(key)
	mp.ct.Add(arr)
}

//countTracer will keep frequency of overall words
type countTracer struct {
	c map[string]int
}

//Add will add a set of modifiers into the count tracer which keeps track of modifier frequency
func (ct *countTracer) Add(w []string) {
	for _, s := range w {
		if l, ok := ct.c[s]; ok {
			l++
			ct.c[s] = l
		} else {
			ct.c[s] = 0
		}
	}
}
