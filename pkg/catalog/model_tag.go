package catalog

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	// EmptyTagError is returned when a raw tag string normalises to nothing.
	EmptyTagError = errors.New("a tag must not be empty")

	tagUnsafeChars = regexp.MustCompile("[^a-z0-9:._-]+")
)

// Tag is a normalised, lowercase, URL-safe label attached to a media. It is identified by
// its normalised string: lowercased, spaces removed, URL-safe ; `namespace:value` is split
// on the first colon ; a tag with no colon is a bare tag (reserved for `favorite`).
type Tag struct {
	namespace string
	value     string
}

// NewTag parses and normalises a raw tag string into a Tag value object. The namespace is
// the part before the first colon ; a raw string without a colon becomes a bare tag whose
// Value is the whole string and whose Namespace is empty.
func NewTag(raw string) (Tag, error) {
	normalised := normaliseTagString(raw)
	if normalised == "" {
		return Tag{}, errors.Wrapf(EmptyTagError, "invalid tag %q", raw)
	}

	if namespace, value, found := strings.Cut(normalised, ":"); found {
		if namespace == "" || value == "" {
			return Tag{}, errors.Wrapf(EmptyTagError, "both namespace and value are mandatory in tag %q", raw)
		}
		return Tag{namespace: namespace, value: value}, nil
	}

	return Tag{value: normalised}, nil
}

// Namespace is the part of the tag before the first colon ; empty for a bare tag.
func (t Tag) Namespace() string {
	return t.namespace
}

// Value is the part of the tag after the first colon ; the whole tag for a bare tag.
func (t Tag) Value() string {
	return t.value
}

// String is the normalised tag, the stable identifier used in storage and URLs.
func (t Tag) String() string {
	if t.namespace == "" {
		return t.value
	}
	return t.namespace + ":" + t.value
}

func normaliseTagString(raw string) string {
	lowered := strings.ToLower(strings.TrimSpace(raw))
	withoutSpaces := strings.ReplaceAll(lowered, " ", "")
	return tagUnsafeChars.ReplaceAllString(withoutSpaces, "")
}
