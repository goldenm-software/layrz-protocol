package layrzprotocol

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func parseArgs(rawArgs string) *map[string]any {
	args := orderedmap.New()

	if rawArgs == "" {
		return &map[string]any{}
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
		case patterns["bleId"].MatchString(key):
			key = strings.Replace(key, ".id", ".mac.address", 1)
		case patterns["bleHumidity"].MatchString(key):
			key = strings.Replace(key, ".hum", ".humidity", 1)
		case patterns["bleTempC"].MatchString(key):
			key = strings.Replace(key, ".tempc", ".temperature.celsius", 1)
		case patterns["bleTempF"].MatchString(key):
			key = strings.Replace(key, ".tempf", ".temperature.fahrenheit", 1)
		case patterns["bleModelId"].MatchString(key):
			key = strings.Replace(key, ".model_id", ".model.id", 1)
		case patterns["bleBatteryLevel"].MatchString(key):
			key = strings.Replace(key, ".batt", ".battery.level", 1)
		case patterns["bleLuxLevel"].MatchString(key):
			key = strings.Replace(key, ".lux", ".light.level.lux", 1)
		case patterns["bleVoltageLevel"].MatchString(key):
			key = strings.Replace(key, ".volt", ".voltage", 1)
		case patterns["bleRpm"].MatchString(key):
			// rpm stays as-is
		case patterns["blePressure"].MatchString(key):
			key = strings.Replace(key, ".press", ".pressure", 1)
		case patterns["bleEventCount"].MatchString(key):
			key = strings.Replace(key, ".counter", ".event.count", 1)
		case patterns["bleXAccel"].MatchString(key):
			key = strings.Replace(key, ".x_acc", ".acceleration.x", 1)
		case patterns["bleYAccel"].MatchString(key):
			key = strings.Replace(key, ".y_acc", ".acceleration.y", 1)
		case patterns["bleZAccel"].MatchString(key):
			key = strings.Replace(key, ".z_acc", ".acceleration.z", 1)
		case patterns["bleMsgCount"].MatchString(key):
			key = strings.Replace(key, ".msg_count", ".message.count", 1)
		case patterns["bleMsg"].MatchString(key):
			key = strings.Replace(key, ".msg", ".message", 1)
		case patterns["bleMagCount"].MatchString(key):
			key = strings.Replace(key, ".mag_counter", ".magnetic.event.count", 1)
		case patterns["bleMagData"].MatchString(key):
			key = strings.Replace(key, ".mag_data", ".magnetic.data", 1)
		case patterns["bleRssi"].MatchString(key):
			key = strings.Replace(key, ".rssi", ".rssi.dbm", 1)
		case key == "report":
			key = "report.code"
		case key == "confiot_ble":
			key = "ble.confiot.connection.status"
		case key == "confiot_serial":
			key = "serial.confiot.connection.status"
		}

		value := strings.Join(subparts[1:], ":")

		intRegexp := regexp.MustCompile(`^-?\d+$`)
		floatRegexp := regexp.MustCompile(`^-?\d+\.\d+$`)

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
	return strings.Replace(strings.Replace(input, "io", "", 1), suffix, "", 1)
}
