# WeatherDate Service
This service is a wrapper to the OpenWeatherAPI (https://home.openweathermap.org). It returns the weather condition for the latitude, longitude point on the globe. The weather condition it returns a map conatining the temperature and conditions such as: `{"city": "Albertville", "outside_condition": "Clouds", "temperature": "31.98 °F", "temperature_condition": "Freezing" }`. It's a simple proof of concept that the a service is made to read weather data and display the weather condition for the (lat,lon) point on the globe.

## OpenWeatherAPI Keys 
API keys are needed for this project. Login or sign up and navigate to the api key page (https://home.openweathermap.org/api_keys). This API key is used in the path parameters of this service. 

## Example of OpenWeatherAPI Call
https://api.openweathermap.org/data/2.5/weather?lat=45.271718&lon=-93.655206&appid=<apiKey>

## Example of weatherWebApp (this project) API Call
http://localhost:8111/v1/weather/data?lat=45.271718&lon=-93.655206&units=imperial&appid=<apiKey>
You can put the API Key in the path parameters directly.

## Run weatherWebApp
Checkout the git repo clesnar/go and rely on the .vscode config - run debug. Or build and run executable.
The default `host:port` is `localhost:8111`.

## debug build
go build -buildvcs=false -o /workspaces/clesnar-go/cmd/weatherWebApp/__debug_bin1800304272 -gcflags all=-N