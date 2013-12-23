package mvc

import (
	"regexp"
	"net/url"
	"strings"
)


type Route struct {
	//host string
	//method  string
    regexp *regexp.Regexp
	keys []string
	controller ControllerInterface
	action string    
}
type StaticDir struct {
	url string
	path string
}
type Router struct {
    routes map[string][]*Route
	staticDirs       map[string][]*StaticDir
}

var (
    // Precompile Regexp to speed things up.
    placeholderMatcher *regexp.Regexp = regexp.MustCompile(`:(\w+)`)
)

func NewRoute(pattern string, controller ControllerInterface, action string) *Route {
	regexp, keys := compilePattern(pattern)
    return &Route{regexp, keys, controller, action}
}

func NewRouter() *Router {
	this := new(Router)
	this.routes = make(map[string][]*Route)
	this.staticDirs = make(map[string][]*StaticDir)
	return this
}

func (this *Router) AddRoute(host string, pattern string, controller ControllerInterface, action string) *Router {
	if host == "" { host = "*" }
	this.routes[host] = append(this.routes[host], NewRoute(pattern, controller, action))
	return this
}

/* add static dir*/
func (this *Router)AddStaticPath(host string, url string, path string) *Router {
	this.staticDirs[host] = append(this.staticDirs[host], &StaticDir{url, path})
	return this
}

func (route *Route) Match(path string) bool {
        return route.regexp.MatchString(path)
}

func (route *Route) extractParams(path string) url.Values {
        params := make(url.Values)
        for i, param := range route.regexp.FindStringSubmatch(path)[1:] {
                params[route.keys[i]] = append(params[route.keys[i]], param)
        }
        return params
}

// compilePattern("/hello/:world") => ^\/hello\/([^#?/]+)$, ["world"]
func compilePattern(pattern string) (*regexp.Regexp, []string) {
        var segments, keys []string
        for _, segment := range strings.Split(pattern, "/") {
                if strings := placeholderMatcher.FindStringSubmatch(segment); strings != nil {
                        keys = append(keys, strings[1])
                        segments = append(segments, placeholderMatcher.ReplaceAllString(segment, "([^#?/]+)"))
                } else {
                        segments = append(segments, segment)
                }
        }
        return regexp.MustCompile(`^` + strings.Join(segments, `\/`) + "$"), keys
}
