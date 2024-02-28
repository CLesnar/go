package model

import (
	"fmt"

	"github.com/CLesnar/go/internal/pkg/util"
)

func (p WeatherConditionGetParameters) Validate() error {
	if p.AppId == nil || len(*p.AppId) < 1 {
		tag := util.FindFieldTagWithDefault(&p, &p.AppId, "json", "appid")
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
