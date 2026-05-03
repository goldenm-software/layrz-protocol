package wire

import (
	"fmt"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		validate func(map[string]any) error
	}{
		{
			name:  "empty args",
			input: "",
			validate: func(args map[string]any) error {
				if len(args) != 0 {
					return fmt.Errorf("expected empty map, got %d items", len(args))
				}
				return nil
			},
		},
		{
			name:  "integer value",
			input: "count:42",
			validate: func(args map[string]any) error {
				val, ok := args["count"]
				if !ok {
					return fmt.Errorf("count key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != 42 {
					return fmt.Errorf("expected int 42, got %v (type %T)", val, val)
				}
				return nil
			},
		},
		{
			name:  "negative integer",
			input: "temperature:-15",
			validate: func(args map[string]any) error {
				val, ok := args["temperature"]
				if !ok {
					return fmt.Errorf("temperature key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != -15 {
					return fmt.Errorf("expected int -15, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "float value",
			input: "latitude:10.123456",
			validate: func(args map[string]any) error {
				val, ok := args["latitude"]
				if !ok {
					return fmt.Errorf("latitude key not found")
				}
				floatVal, ok := val.(float64)
				if !ok || floatVal != 10.123456 {
					return fmt.Errorf("expected float64 10.123456, got %v (type %T)", val, val)
				}
				return nil
			},
		},
		{
			name:  "boolean true",
			input: "enabled:true",
			validate: func(args map[string]any) error {
				val, ok := args["enabled"]
				if !ok {
					return fmt.Errorf("enabled key not found")
				}
				boolVal, ok := val.(bool)
				if !ok || !boolVal {
					return fmt.Errorf("expected bool true, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "boolean false",
			input: "disabled:false",
			validate: func(args map[string]any) error {
				val, ok := args["disabled"]
				if !ok {
					return fmt.Errorf("disabled key not found")
				}
				boolVal, ok := val.(bool)
				if !ok || boolVal {
					return fmt.Errorf("expected bool false, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "string value",
			input: "name:hello",
			validate: func(args map[string]any) error {
				val, ok := args["name"]
				if !ok {
					return fmt.Errorf("name key not found")
				}
				strVal, ok := val.(string)
				if !ok || strVal != "hello" {
					return fmt.Errorf("expected string hello, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "GPIO digital input remapping",
			input: "io3.di:1",
			validate: func(args map[string]any) error {
				val, ok := args["gpio.3.digital.input"]
				if !ok {
					return fmt.Errorf("gpio.3.digital.input key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != 1 {
					return fmt.Errorf("expected int 1, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "GPIO digital output remapping",
			input: "io5.do:0",
			validate: func(args map[string]any) error {
				_, ok := args["gpio.5.digital.output"]
				if !ok {
					return fmt.Errorf("gpio.5.digital.output key not found")
				}
				return nil
			},
		},
		{
			name:  "GPIO analog input remapping",
			input: "io2.ai:512",
			validate: func(args map[string]any) error {
				_, ok := args["gpio.2.analog.input"]
				if !ok {
					return fmt.Errorf("gpio.2.analog.input key not found")
				}
				return nil
			},
		},
		{
			name:  "GPIO counter remapping",
			input: "io1.counter:999",
			validate: func(args map[string]any) error {
				_, ok := args["gpio.1.event.count"]
				if !ok {
					return fmt.Errorf("gpio.1.event.count key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE temperature celsius remapping",
			input: "ble.0.tempc:25",
			validate: func(args map[string]any) error {
				val, ok := args["ble.0.temperature.celsius"]
				if !ok {
					return fmt.Errorf("ble.0.temperature.celsius key not found")
				}
				intVal, ok := val.(int)
				if !ok || intVal != 25 {
					return fmt.Errorf("expected int 25, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "BLE humidity remapping",
			input: "ble.1.hum:60.5",
			validate: func(args map[string]any) error {
				_, ok := args["ble.1.humidity"]
				if !ok {
					return fmt.Errorf("ble.1.humidity key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE id (mac address) remapping",
			input: "ble.2.id:AA:BB:CC:DD:EE:FF",
			validate: func(args map[string]any) error {
				val, ok := args["ble.2.mac.address"]
				if !ok {
					return fmt.Errorf("ble.2.mac.address key not found")
				}
				strVal, ok := val.(string)
				if !ok || strVal != "AA:BB:CC:DD:EE:FF" {
					return fmt.Errorf("expected string, got %v", val)
				}
				return nil
			},
		},
		{
			name:  "multiple args",
			input: "count:42,enabled:true,name:test,temperature:23.5",
			validate: func(args map[string]any) error {
				if len(args) != 4 {
					return fmt.Errorf("expected 4 args, got %d", len(args))
				}
				return nil
			},
		},
		{
			name:  "report code remapping",
			input: "report:123",
			validate: func(args map[string]any) error {
				_, ok := args["report.code"]
				if !ok {
					return fmt.Errorf("report.code key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE temperature fahrenheit remapping",
			input: "ble.0.tempf:77",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.temperature.fahrenheit"]
				if !ok {
					return fmt.Errorf("ble.0.temperature.fahrenheit key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE model_id remapping",
			input: "ble.0.model_id:42",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.model.id"]
				if !ok {
					return fmt.Errorf("ble.0.model.id key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE battery remapping",
			input: "ble.0.batt:85",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.battery.level"]
				if !ok {
					return fmt.Errorf("ble.0.battery.level key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE lux remapping",
			input: "ble.0.lux:500",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.light.level.lux"]
				if !ok {
					return fmt.Errorf("ble.0.light.level.lux key not found")
				}
				return nil
			},
		},
		{
			name:  "GPIO analog output remapping",
			input: "io4.ao:2048",
			validate: func(args map[string]any) error {
				_, ok := args["gpio.4.analog.output"]
				if !ok {
					return fmt.Errorf("gpio.4.analog.output key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE voltage level remapping",
			input: "ble.0.volt:3300",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.voltage"]
				if !ok {
					return fmt.Errorf("ble.0.voltage key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE rpm remapping (key stays)",
			input: "ble.0.rpm:1200",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.rpm"]
				if !ok {
					return fmt.Errorf("ble.0.rpm key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE pressure remapping",
			input: "ble.1.press:1013",
			validate: func(args map[string]any) error {
				_, ok := args["ble.1.pressure"]
				if !ok {
					return fmt.Errorf("ble.1.pressure key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE event count remapping",
			input: "ble.2.counter:5",
			validate: func(args map[string]any) error {
				_, ok := args["ble.2.event.count"]
				if !ok {
					return fmt.Errorf("ble.2.event.count key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE X acceleration remapping",
			input: "ble.0.x_acc:100",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.acceleration.x"]
				if !ok {
					return fmt.Errorf("ble.0.acceleration.x key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE Y acceleration remapping",
			input: "ble.0.y_acc:-50",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.acceleration.y"]
				if !ok {
					return fmt.Errorf("ble.0.acceleration.y key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE Z acceleration remapping",
			input: "ble.0.z_acc:200",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.acceleration.z"]
				if !ok {
					return fmt.Errorf("ble.0.acceleration.z key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE message count remapping",
			input: "ble.0.msg_count:3",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.message.count"]
				if !ok {
					return fmt.Errorf("ble.0.message.count key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE message remapping",
			input: "ble.0.msg:hello",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.message"]
				if !ok {
					return fmt.Errorf("ble.0.message key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE magnetic event count remapping",
			input: "ble.0.mag_counter:7",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.magnetic.event.count"]
				if !ok {
					return fmt.Errorf("ble.0.magnetic.event.count key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE magnetic data remapping",
			input: "ble.0.mag_data:abc123",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.magnetic.data"]
				if !ok {
					return fmt.Errorf("ble.0.magnetic.data key not found")
				}
				return nil
			},
		},
		{
			name:  "BLE RSSI remapping",
			input: "ble.0.rssi:-70",
			validate: func(args map[string]any) error {
				_, ok := args["ble.0.rssi.dbm"]
				if !ok {
					return fmt.Errorf("ble.0.rssi.dbm key not found")
				}
				return nil
			},
		},
		{
			name:  "confiot_ble remapping",
			input: "confiot_ble:1",
			validate: func(args map[string]any) error {
				_, ok := args["ble.confiot.connection.status"]
				if !ok {
					return fmt.Errorf("ble.confiot.connection.status key not found")
				}
				return nil
			},
		},
		{
			name:  "confiot_serial remapping",
			input: "confiot_serial:1",
			validate: func(args map[string]any) error {
				_, ok := args["serial.confiot.connection.status"]
				if !ok {
					return fmt.Errorf("serial.confiot.connection.status key not found")
				}
				return nil
			},
		},
		{
			name:  "malformed token without colon is skipped",
			input: "keyonly,valid:42",
			validate: func(args map[string]any) error {
				if _, ok := args["keyonly"]; ok {
					return fmt.Errorf("keyonly should have been skipped")
				}
				if _, ok := args["valid"]; !ok {
					return fmt.Errorf("valid key should be present")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseArgs(tt.input)
			if err := tt.validate(result); err != nil {
				t.Errorf("validation failed: %v", err)
			}
		})
	}
}
