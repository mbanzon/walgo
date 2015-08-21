package walgo

import (
	"net/url"
)

func createParameterUrl(urlStr string, p ParameterMap) (u *url.URL, err error) {
	u, err = url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	query := u.Query()

	if p != nil {
		for k, v := range p {
			query.Add(k, v)
		}
	}

	u.RawQuery = query.Encode()
	return u, nil
}
