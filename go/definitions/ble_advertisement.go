package definitions

import "time"

// BleAdvertisement defines a detection of a BLE advertisement
type BleAdvertisement struct {
	// Is the detected Mac Address. This Mac Adress came from the detected device, not the
	// device that is detecting.
	MacAddress string
	// Is when the device was detected.
	Timestamp time.Time

	// Is the closest latitude of the device. Defined by the device that is detecting.
	// This value is optional
	Latitude *float64

	// Is the closest longitude of the device. Defined by the device that is detecting.
	// This value is optional
	Longitude *float64

	// Is the closest altitude of the device. Defined by the device that is detecting.
	// This value is optional
	Altitude *float64

	// Is the signal strength of the detected device.
	Rssi int

	// Is the transmission power of the detected device.
	// This value is optional
	TxPower int

	// Is the model of the detected device. This model should be equals to the model of the device
	// and the model defined by Layrz.
	Model string

	// Is the list of manufacturer data advertised by the device.
	ManufacturerData []BleManufacturerData

	// Is the list of service data advertised by the device.
	ServiceData []BleServiceData

	// Is the name of the device. This name is optional.
	DeviceName string
}

// BleManufacturerData defines the manufacturer data advertised by the device
type BleManufacturerData struct {
	// Is the manufacturer identifier.
	CompanyId int

	// Is the manufacturer data.
	Data []byte
}

// BleServiceData defines the service data advertised by the device
type BleServiceData struct {
	// Is the service UUID.
	Uuid int

	// Is the service data.
	Data []byte
}
