package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Route represents a route pattern with HTTP method
type Route struct {
	Pattern string // Pattern like "/api/v1/owners/{owner}/albums/{folderName}"
	Method  string // HTTP method like "GET", "POST", etc.
}

// MatchedRoute represents a matched route with extracted path parameters
type MatchedRoute struct {
	Route      Route
	PathParams map[string]string
}

// routePattern holds a compiled route pattern
type routePattern struct {
	route       Route
	regex       *regexp.Regexp
	paramNames  []string
}

// compileRoute compiles a route pattern into a regex and extracts parameter names
func compileRoute(route Route) (*routePattern, error) {
	// Extract parameter names from the pattern
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(route.Pattern, -1)
	
	paramNames := make([]string, 0, len(matches))
	for _, match := range matches {
		paramNames = append(paramNames, match[1])
	}
	
	// Convert the pattern to a regex pattern
	// Replace {param} with a capturing group
	regexPattern := regexp.QuoteMeta(route.Pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\{`, "{")
	regexPattern = strings.ReplaceAll(regexPattern, `\}`, "}")
	regexPattern = paramRegex.ReplaceAllString(regexPattern, `([^/]+)`)
	regexPattern = "^" + regexPattern + "$"
	
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile route pattern %s: %w", route.Pattern, err)
	}
	
	return &routePattern{
		route:      route,
		regex:      regex,
		paramNames: paramNames,
	}, nil
}

// MatchRoute attempts to match a method and path against a list of routes
// Returns the matched route with extracted path parameters, or an error if no match is found
func MatchRoute(routes []Route, method, path string) (*MatchedRoute, error) {
	// Compile all routes
	patterns := make([]*routePattern, 0, len(routes))
	for _, route := range routes {
		pattern, err := compileRoute(route)
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, pattern)
	}
	
	// Try to match each pattern
	for _, pattern := range patterns {
		// Check if method matches
		if pattern.route.Method != method {
			continue
		}
		
		// Check if path matches
		matches := pattern.regex.FindStringSubmatch(path)
		if matches == nil {
			continue
		}
		
		// Extract path parameters
		pathParams := make(map[string]string)
		for i, paramName := range pattern.paramNames {
			pathParams[paramName] = matches[i+1] // Skip the first match (full string)
		}
		
		return &MatchedRoute{
			Route:      pattern.route,
			PathParams: pathParams,
		}, nil
	}
	
	return nil, fmt.Errorf("no matching route found for %s %s", method, path)
}
