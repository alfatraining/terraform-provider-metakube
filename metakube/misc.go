package metakube

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
)

func labelsMap(d *schema.ResourceData) (ret map[string]string) {
	if attr, ok := d.GetOk("labels"); ok {
		ret = make(map[string]string)
		for k, v := range attr.(map[string]interface{}) {
			ret[k] = v.(string)
		}
	}
	return ret
}

func clusterVersionsHasPrefix(version, prefix string) bool {
	return len(version) >= len(prefix) && version[:len(prefix)] == prefix
}

type clusterVersionsParsed struct {
	p0 int64
	p1 int64
	p2 int64
}

func parseClusterVersion(version string) (*clusterVersionsParsed, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, errors.Errorf("unknown version format: %v", version)
	}
	parsed := make([]int64, 0)
	for _, p := range parts {
		v, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "parse version `%s`", version)
		}
		parsed = append(parsed, v)
	}
	return &clusterVersionsParsed{
		p0: parsed[0],
		p1: parsed[1],
		p2: parsed[2],
	}, nil
}

func clusterVersionBigger(a, b string) (bool, error) {
	if aParsed, err := parseClusterVersion(a); err != nil {
		return false, err
	} else if bParsed, err := parseClusterVersion(b); err != nil {
		return false, err
	} else if aParsed.p0 > bParsed.p0 {
		return true, nil
	} else if aParsed.p0 == bParsed.p0 && aParsed.p1 > bParsed.p1 {
		return true, nil
	} else if aParsed.p0 == bParsed.p0 && aParsed.p1 == bParsed.p1 && aParsed.p2 > bParsed.p2 {
		return true, nil
	} else {
		return false, nil
	}
}
