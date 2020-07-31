package mws

import (
	"strings"

	"github.com/pkg/errors"
)

type region struct {
	regionID      string
	country       string
	endpoint      string
	marketPlaceID string
}

func findRegionByCountry(country string) (r region, err error) {
	for _, region := range regions {
		if strings.EqualFold(region.country, country) {
			r = region
			return
		}
	}
	err = errors.Errorf("region not found ,invalid country = %s", country)
	return
}

var regions = []region{
	{"NA", "US", "https://mws.amazonservices.com/", "ATVPDKIKX0DER"},
	{"NA", "CA", "https://mws.amazonservices.ca/", "A2EUQ1WTGCTBG2"},
	{"NA", "MX", "https://mws.amazonservices.com.mx/", "A1AM78C64UM0Y8"},
	{"NA", "BR", "https://mws.amazonservices.com/", "A2Q3Y263D00KWC"},

	{"EU", "AE", "https://mws.amazonservices.ae/", "A2VIGQ35RCS4UG"},
	{"EU", "DE", "https://mws-eu.amazonservices.com/", "A1PA6795UKMFR9"},
	{"EU", "EG", "https://mws-eu.amazonservices.com/", "ARBP9OOSHTCHU"},
	{"EU", "ES", "https://mws-eu.amazonservices.com/", "A1RKKUPIHCS9HS"},
	{"EU", "FR", "https://mws-eu.amazonservices.com/", "A13V1IB3VIYZZH"},
	{"EU", "GB", "https://mws-eu.amazonservices.com/", "A1F83G8C2ARO7P"},
	{"EU", "IN", "https://mws.amazonservices.in/", "A21TJRUUN4KGV"},
	{"EU", "IT", "https://mws-eu.amazonservices.com/", "APJ6JRA9NG5V4"},
	{"EU", "NL", "https://mws-eu.amazonservices.com/", "A1805IZSGTT6HS"},
	{"EU", "SA", "https://mws-eu.amazonservices.com/", "A17E79C6D8DWNP"},
	{"EU", "TR", "https://mws-eu.amazonservices.com/", "A33AVAJ2PDY3EV"},

	{"EF", "JP", "https://mws.amazonservices.jp/", "A1VC38T7YXB528"},
	{"EF", "AU", "https://mws.amazonservices.com.au/", "A39IBJ37TRP1C6"},
	{"EF", "SG", "	https://mws-fe.amazonservices.com/", "A19VAU5U5O7RUS"},
}
