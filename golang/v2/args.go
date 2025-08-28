package layrzprotocol

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func parseArgs(rawArgs string) *map[string]interface{} {
	args := orderedmap.New()

	if rawArgs == "" {
		return &map[string]interface{}{}
	}

	parts := strings.Split(rawArgs, ",")
	for _, part := range parts {
		subparts := strings.Split(part, ":")
		if len(subparts) < 2 {
			continue
		}

		key := subparts[0]

		patterns := map[string]*regexp.Regexp{
			"digitalInput":    regexp.MustCompile(`^io[0-9]+\.di$`),
			"digitalOutput":   regexp.MustCompile(`^io[0-9]+\.do$`),
			"analogInput":     regexp.MustCompile(`^io[0-9]+\.ai$`),
			"analogOutput":    regexp.MustCompile(`^io[0-9]+\.ao$`),
			"counter":         regexp.MustCompile(`^io[0-9]+\.counter$`),
			"bleId":           regexp.MustCompile(`^ble\.[0-9]+\.id$`),
			"bleHumidity":     regexp.MustCompile(`^ble\.[0-9]+\.hum$`),
			"bleTempC":        regexp.MustCompile(`^ble\.[0-9]+\.tempc$`),
			"bleTempF":        regexp.MustCompile(`^ble\.[0-9]+\.tempf$`),
			"bleModelId":      regexp.MustCompile(`^ble\.[0-9]+\.model_id$`),
			"bleBatteryLevel": regexp.MustCompile(`^ble\.[0-9]+\.batt$`),
			"bleLuxLevel":     regexp.MustCompile(`^ble\.[0-9]+\.lux$`),
			"bleVoltageLevel": regexp.MustCompile(`^ble\.[0-9]+\.volt$`),
			"bleRpm":          regexp.MustCompile(`^ble\.[0-9]+\.rpm$`),
			"blePressure":     regexp.MustCompile(`^ble\.[0-9]+\.press$`),
			"bleEventCount":   regexp.MustCompile(`^ble\.[0-9]+\.counter$`),
			"bleXAccel":       regexp.MustCompile(`^ble\.[0-9]+\.x_acc$`),
			"bleYAccel":       regexp.MustCompile(`^ble\.[0-9]+\.y_acc$`),
			"bleZAccel":       regexp.MustCompile(`^ble\.[0-9]+\.z_acc$`),
			"bleMsgCount":     regexp.MustCompile(`^ble\.[0-9]+\.msg_count$`),
			"bleMsg":          regexp.MustCompile(`^ble\.[0-9]+\.msg$`),
			"bleMagCount":     regexp.MustCompile(`^ble\.[0-9]+\.mag_counter$`),
			"bleMagData":      regexp.MustCompile(`^ble\.[0-9]+\.mag_data$`),
			"bleRssi":         regexp.MustCompile(`^ble\.[0-9]+\.rssi$`),
		}

		switch {
		case patterns["digitalInput"].MatchString(key):
			key = "gpio." + extractGpio(key, ".di") + ".digital.input"
		case patterns["digitalOutput"].MatchString(key):
			key = "gpio." + extractGpio(key, ".do") + ".digital.output"
		case patterns["analogInput"].MatchString(key):
			key = "gpio." + extractGpio(key, ".ai") + ".analog.input"
		case patterns["analogOutput"].MatchString(key):
			key = "gpio." + extractGpio(key, ".ao") + ".analog.output"
		case patterns["counter"].MatchString(key):
			key = "gpio." + extractGpio(key, ".counter") + ".event.count"
		case key == "report":
			key = "report.code"
		case key == "confiot_ble":
			key = "ble.confiot.connection.status"
		case key == "confiot_serial":
			key = "serial.confiot.connection.status"
		}

		value := strings.Join(subparts[1:], ":")

		intRegexp := regexp.MustCompile(`^\d+$`)
		floatRegexp := regexp.MustCompile(`^\d+\.\d+$`)

		if intRegexp.MatchString(value) {
			intVal, err := strconv.Atoi(value)
			if err == nil {
				args.Set(key, intVal)
			}
		} else if floatRegexp.MatchString(value) {
			floatVal, err := strconv.ParseFloat(value, 64)
			if err == nil {
				args.Set(key, floatVal)
			}
		} else if (value == "true") || (value == "false") {
			args.Set(key, value == "true")
		} else {
			args.Set(key, value)
		}
	}

	output := args.Values()
	return &output
}

func extractGpio(input, suffix string) string {
	return strings.Replace(strings.Replace(input, "io.", "", 1), suffix, "", 1)
}
