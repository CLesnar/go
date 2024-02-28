package weather

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CLesnar/go/internal/pkg/http"
	"github.com/CLesnar/go/internal/pkg/util"
	"github.com/CLesnar/go/pkg/model"
)

// openweatherapi is one source for getting weather data.

// See OpenWeatherMap details on APIs: https://openweathermap.org/current
// Example of API call to CurrentWeatherData:
// https://api.openweathermap.org/data/2.5/weather?lat=44.34&lon=10.99&appid={API key}

// implements weather interface
type OpenWeatherMap struct {
}

func (w *OpenWeatherMap) GetCurrentWeatherData(ctx context.Context, parameters model.OpenWeatherMapParametersGetCurrentData) (model.OpenWeatherMapResponse, error) {
	httpHelper, respData := http.Http{}, model.OpenWeatherMapResponse{}
	resp, err := httpHelper.GetRequest(ctx, OpenWeatherMapBaseUrl, nil, parameters, nil)
	if err != nil {
		return respData, err
	}
	if err := json.Unmarshal(resp, &respData); err != nil {
		return respData, err
	}
	return respData, nil
}

func (w *OpenWeatherMap) GetWeatherCondition(ctx context.Context, lat float64, lon float64, apiKey string) (map[string]interface{}, error) {
	parameters := model.OpenWeatherMapParametersGetCurrentData{AppId: util.ConstantRef(apiKey), Latitude: util.ConstantRef(lat), Longitude: util.ConstantRef(lon), Units: util.ConstantRef("imperial")}
	weatherConditionMap := map[string]interface{}{}
	data, err := w.GetCurrentWeatherData(ctx, parameters)
	if err != nil {
		return weatherConditionMap, err
	}
	weatherConditionMap["temperature"] = fmt.Sprintf("%v °F", data.Main.Temperature)
	weatherConditionMap["temperature_condition"] = temperatureCondition(data.Main.FeelsLike)
	weatherConditionMap["city"] = data.Name
	outsideCondition := ""
	for _, weather := range data.Weather {
		outsideCondition += ", " + weather.Main
	}
	if len(outsideCondition) <= 2 {
		outsideCondition = "Normal"
	} else {
		outsideCondition = outsideCondition[2:] // remove first ", "
	}
	weatherConditionMap["outside_condition"] = outsideCondition
	return weatherConditionMap, nil
}
