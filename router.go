package zapi

import (
	"regexp"
	"strings"
)

type route struct {
	method     string
	regex      *regexp.Regexp
	handler    HandlerFunc
	paramNames []string
}

type router struct {
	routes []route
}

func (r *router) add(method, path string, h HandlerFunc) {

	var paramNames []string
	pattern := "^"
	segments := strings.Split(strings.Trim(path, "/"), "/")

	// Track if this route is more specific than existing ones
	isMoreSpecific := false

	for _, segment := range segments {
		pattern += "/"

		// Check for parameter patterns like @[id] or [id]
		if paramMatch := regexp.MustCompile(`^(.*)\[([^\]]+)\](.*)$`).FindStringSubmatch(segment); paramMatch != nil {
			prefix := regexp.QuoteMeta(paramMatch[1])
			paramName := paramMatch[2]
			suffix := regexp.QuoteMeta(paramMatch[3])

			paramNames = append(paramNames, paramName)
			pattern += prefix + `([^/]+)` + suffix

			// If this segment has prefixes/suffixes, it's more specific
			if prefix != "" || suffix != "" {
				isMoreSpecific = true
			}
		} else if segment == "*" {
			// Handle wildcard (least specific)
			paramNames = append(paramNames, "*")
			pattern += `(.+)`
		} else {
			// Regular path segment (most specific)
			pattern += regexp.QuoteMeta(segment)
			isMoreSpecific = true
		}
	}

	pattern += "/?$"

	entry := route{
		method:     method,
		regex:      regexp.MustCompile(pattern),
		handler:    h,
		paramNames: paramNames,
	}

	if isMoreSpecific {
		r.routes = append([]route{entry}, r.routes...)
	} else {
		r.routes = append(r.routes, entry)
	}
}

func (r *router) get(path string, h HandlerFunc) {
	r.add("GET", path, h)
}

func (r *router) post(path string, h HandlerFunc) {
	r.add("POST", path, h)
}

func (r *router) put(path string, h HandlerFunc) {
	r.add("PUT", path, h)
}

func (r *router) delete(path string, h HandlerFunc) {
	r.add("DELETE", path, h)
}
