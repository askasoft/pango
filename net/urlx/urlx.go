package urlx

import (
	"net/url"
)

// CleanURL returns the URL without user, query and fragment.
func CleanURL(link string) (string, error) {
	ulink, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	ulink.ForceQuery = false
	ulink.RawQuery = ""
	ulink.Fragment = ""
	ulink.User = nil

	return ulink.String(), nil
}

// ParentURL returns the parent URL relative to the `link`.
func ParentURL(link string) (string, error) {
	ulink, err := url.Parse(link)
	if err != nil {
		return link, err
	}

	ulink.ForceQuery = false
	ulink.RawQuery = ""
	ulink.Fragment = ""
	ulink.User = nil
	if ulink.Path == "" {
		ulink.Path = "/"
		return ulink.String(), nil
	}

	udir, err := url.Parse("./")
	if err != nil {
		return link, err
	}

	ubase := ulink.ResolveReference(udir)
	return ubase.String(), nil
}
