package model

import (
	"fmt"

	"github.com/CLesnar/go/internal/pkg/util"
)

func (p WeatherConditionParameters) Validate() error {
	if p.ApiId == nil || len(*p.ApiId) < 1 {
		tag := util.FindFieldTagWithDefault(&p, &p.ApiId, "json", "appid")
		return fmt.Errorf("missing required path parameter '%s'", tag)
	}
	if p.Latitude == nil {
		tag := util.FindFieldTagWithDefault(&p, &p.Latitude, "json", "lat")
		return fmt.Errorf("missing required path parameter '%s'", tag)
	}
	if p.Longitude == nil {
		tag := util.FindFieldTagWithDefault(&p, &p.Longitude, "json", "lon")
		return fmt.Errorf("missing required path parameter '%s'", tag)
	}
	return nil
}
