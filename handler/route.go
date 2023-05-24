package handler

import (
	"log"
	"net/http"
	"regexp"
	"sort"
	"fmt"
)

type RouteHandlerFunc func() (http.Handler, error)

func RouteHandler(handlers map[string]RouteHandlerFunc) (http.Handler, error) {

	patterns := make([]string, 0)
	re_lookup := make(map[string]*regexp.Regexp)
	
	for re_pat, _ := range handlers {

		re, err := regexp.Compile(re_pat)

		if err != nil {
			return nil, fmt.Errorf("Failed to compile '%s', %w", re_pat, err)
		}

		re_lookup[re_pat] = re
		patterns = append(patterns, re_pat)
	}

	sort.Strings(patterns)
	
	fn := func(rsp http.ResponseWriter, req *http.Request) {

		var handler_func RouteHandlerFunc

		path := req.URL.Path

		for _, re_pat := range patterns {

			re, ok := re_lookup[re_pat]

			if !ok {
				continue
			}
			
			if !re.MatchString(path) {
				continue
			}

			handler_func = handlers[re_pat]
			break
		}

		if handler_func == nil {
			http.Error(rsp, "Not found", http.StatusNotFound)
			return
		}

		h, err := handler_func()

		if err != nil {
			log.Printf("Failed to instantiate handler func for '%s', %v", path, err)
			http.Error(rsp, "Internal server error", http.StatusInternalServerError)
			return
		}

		h.ServeHTTP(rsp, req)
		return
	}

	return http.HandlerFunc(fn), nil
}
