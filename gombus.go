package gombus

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	OTHER              = "Other"
	OIL                = "Oil"
	ELECTRICITY        = "Electricity"
	GAS                = "Gas"
	HEAT_OUTLET        = "Heat (Volume measured at return temperature: outlet)"
	STREAM             = "Steam"
	HOT_WATER          = "Hot Water"
	WATER              = "Water"
	HEAT_COST_ALLOC    = "Heat Cost Allocator."
	COMPRESSED_AIR     = "Compressed Air"
	COLLING_RET_OUTLET = "Cooling load meter (Volume measured at return temperature: outlet)"
	COLLING_FLOW_INLET = "Cooling load meter (Volume measured at flow temperature: inlet)"
	HEAT_INLET         = "Heat (Volume measured at flow temperature: inlet)"
	HEAT_COOLING       = "Heat / Cooling load meter"
	BUS_SYSTEM         = "Bus / System"
	UNKNOWN_MEDIUM     = "Unknown Medium"
	COLD_WATER         = "Cold Water"
	DUAL_WATER         = "Dual Water"
	PRESSURE           = "Pressure"
	A_D_CONVERTER      = "A/D Converter"
	RESERVED           = "Reserved"
)

// VIF codes
const (
	PARAMETER_UNDEFINED          = "PARAMETER_UNDEFINED"
	PARAMETER_ENERGY             = "PARAMETER_ENERGY"
	PARAMETER_MASS               = "PARAMETER_MASS"
	PARAMETER_POWER              = "PARAMETER_POWER"
	PARAMETER_VOLUME             = "PARAMETER_VOLUME"
	PARAMETER_VOLUME_FLOW        = "PARAMETER_VOLUME_FLOW"
	PARAMETER_MASS_FLOW          = "PARAMETER_MASS_FLOW"
	PARAMETER_TEMP_FLOW          = "PARAMETER_TEMP_FLOW"
	PARAMETER_TEMP_RETURN        = "PARAMETER_TEMP_RETURN"
	PARAMETER_PRESSURE           = "PARAMETER_PRESSURE"
	PARAMETER_ON_TIME            = "PARAMETER_ON_TIME"
	PARAMETER_OPERATING_TIME     = "PARAMETER_OPERATING_TIME"
	PARAMETER_AVERAGING_DURATION = "PARAMETER_AVERAGING_DURATION"
	PARAMETER_ACTUALITY_DURATION = "PARAMETER_ACTUALITY_DURATION"
	PARAMETER_DATETIME           = "PARAMETER_DATETIME"
	PARAMETER_DATE               = "PARAMETER_DATE"
	PARAMETER_TEMP_DIFF          = "PARAMETER_TEMP_DIFF"
	PARAMETER_TEMP_EXTERNAL      = "PARAMETER_TEMP_EXTERNAL"
	PARAMETER_UNITS              = "PARAMETER_UNITS"
	PARAMETER_RESERVED           = "PARAMETER_RESERVED"
	PARAMETER_CUSTOM_VIF         = "PARAMETER_CUSTOM_VIF"
	PARAMETER_FABRICATION        = "PARAMETER_FABRICATION"
	PARAMETER_BUS_ADDR           = "PARAMETER_BUS_ADDR"
	PARAMETER_MANUFACTURED_SPEC  = "PARAMETER_MANUFACTURED_SPEC"
	PARAMETER_FINDWARE_VERSION   = "PARAMETER_FINDWARE_VERSION"
	PARAMETER_SOFTWARE_VERSION   = "PARAMETER_SOFTWARE_VERSION"
	PARAMETER_ACCESS_NUMBER      = "PARAMETER_ACCESS_NUMBER"
	PARAMETER_MEDIUM             = "PARAMETER_MEDIUM"
	PARAMETER_MANUFACTURER       = "PARAMETER_MANUFACTURER"
	PARAMETER_IS_IDENTIFICATION  = "PARAMETER_IS_IDENTIFICATION"
	PARAMETER_MODEL_VERSION      = "PARAMETER_MODEL_VERSION"
	PARAMETER_HARDWARE_VERSION   = "PARAMETER_HARDWARE_VERSION"
	PARAMETER_PASSWORD           = "PARAMETER_PASSWORD"
	PARAMETER_ERROR_FLAG         = "PARAMETER_ERROR_FLAG"
	PARAMETER_CUSTOMER_LOCATION  = "PARAMETER_CUSTOMER_LOCATION"
	PARAMETER_CUSTOMER           = "PARAMETER_CUSTOMER"
	PARAMETER_DIGITAL_OUTPUT     = "PARAMETER_DIGITAL_OUTPUT"
	PARAMETER_DIGITAL_INPUT      = "PARAMETER_DIGITAL_INPUT"
	PARAMETER_V                  = "PARAMETER_V"
	PARAMETER_A                  = "PARAMETER_A"
	PARAMETER_UNRECOGNIZED       = "PARAMETER_UNRECOGNIZED"
)

/*
	Units
*/
const (
	UNIT_HMS = "h,m,s"
	UNIT_DMY = "D,M,Y"
	UNIT_WH  = "Wh"
	UNIT_KWH = "kWh"
	UNIT_MWH = "MWh"
	UNIT_KJ  = "kJ"
	UNIT_MJ  = "MJ"
	UNIT_GJ  = "GJ"
	UNIT_W   = "W"
	UNIT_KW  = "kW"
	UNIT_MW  = "MW"
	UNIT_KJH = "kJ/h"
)

const (
	EXP_NONE = 0
	EXP_m    = 1
	EXP_MY   = 2
	EXP_10   = 3
	EXP_100  = 4
	EXP_K    = 5
	EXP_10K  = 6
	EXP_100K = 7
	EXP_M    = 8
	EXP_T    = 9
	EXP_1E   = 10
)

type DataRecord struct {
	parameter  string
	value      string
	conversion string
	unit       string
}

type SlaveInformation struct {
	id           int
	manufactured string
	version      string
	product_name string
	medium       string
	accessnumber int
	status       int
	signature    int
}

type Mbus struct {
	SlaveInformation
	data   []*DataRecord
	num485 string
}

func SplitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}

//
//-----------------------
// Generate request
//-----------------------
//
func (m *Mbus) GetPackConnect() []byte {
	res := []byte{0x68, 0x0B, 0x0B, 0x68, 0x53, 0xFD, 0x52}

	v := fmt.Sprintf("%08s", m.num485)

	r := SplitSubN(v, 2)

	for i := 4; i > 0; i-- {
		var b byte
		fmt.Sscanf(r[i-1], "%X", &b)
		res = append(res, b)
	}

	res = append(res, 0xFF)
	res = append(res, 0xFF)
	res = append(res, 0xFF)
	res = append(res, 0xFF)
	res = append(res, byte(crc(res[4:])&0xFF))
	res = append(res, 0x16)
	return res
}

func (m *Mbus) GetPackReadData() []byte {
	res := []byte{}

	var code uint8 = 0x7B // 0x7B
	code_num := 0xFD

	crc := int(code) + int(code_num)
	res = append(res, 0x10)
	res = append(res, code)
	res = append(res, byte(code_num))
	res = append(res, byte(crc&0xFF))
	res = append(res, 0x16)

	return res
}

//
//-----------------------
// Mbus get functions
//-----------------------
//

func (m *Mbus) SetNum485(num485 string) {
	m.num485 = num485
}

func (m *Mbus) GetId() int {
	return m.SlaveInformation.id
}

func (m *Mbus) GetManufacturer() string {
	return m.SlaveInformation.manufactured
}

func (m *Mbus) GetVersion() string {
	return m.SlaveInformation.version
}

func (m *Mbus) GetMedium() string {
	return m.SlaveInformation.medium
}

func (m *Mbus) GetProductName() string {
	return m.SlaveInformation.product_name
}

func (m *Mbus) GetAccessNumber() int {
	return m.SlaveInformation.accessnumber
}

func (m *Mbus) GetSignature() int {
	return m.SlaveInformation.signature
}

func (m *Mbus) GetData() []*DataRecord {
	return m.data
}

//
//-----------------------
// Data record get functions
//-----------------------
//

func (d *DataRecord) GetParameterIdent() string {
	return d.parameter
}

func (d *DataRecord) GetValue() string {
	return d.value
}

func (d *DataRecord) GetUnit() string {
	return d.unit
}

func (d *DataRecord) GetConversion() string {
	return d.conversion
}

func New() *Mbus {
	return &Mbus{}
}

func (m *Mbus) ParseFrames(answer []byte) {

	if len(answer) < 19 {
		return
	}

	if IsValidCRC(answer) == false {
		return
	}

	if answer[0] != 0x68 || answer[3] != 0x68 {
		return
	}

	if answer[1] != answer[2] {
		return
	}

	if len(answer[5:len(answer)-1]) != int(answer[1]) {
		return
	}

	m.ParseHead(answer)
	m.ParseDataRecords(answer[19:])

}

func IsValidCRC(answer []byte) bool {

	if len(answer) == 1 {
		if answer[0] == 0xE5 {
			return true
		} else {
			return false
		}
	}

	if len(answer) < 2 {
		return false
	}

	if len(answer) != int(answer[1])+6 /* 68 + len + len + 68 + len DATA + CRC + 16 */ {
		return false
	}

	res := 0
	data := answer[4 : len(answer)-2]
	for i := 0; i < len(data); i++ {
		res += int(data[i])
	}

	return answer[len(answer)-2] == byte(res&0xFF)
}

func crc(cmd []byte) int {
	res := 0
	for i := 0; i < len(cmd); i++ {
		res += int(cmd[i])
	}
	return res
}

func (m *Mbus) ParseHead(answer []byte) {
	str := fmt.Sprintf("%02X%02X%02X%02X", answer[10], answer[9], answer[8], answer[7])
	i, err := strconv.Atoi(str)
	if err != nil {
		return
	}
	m.SlaveInformation.id = i

	m.manufactured = m.mbus_decode_manufacturer(answer[11], answer[12])

	m.version = string(answer[13])
	m.medium = func(c byte) string {

		switch c {
		case 0x00:
			return OTHER
		case 0x01:
			return OIL
		case 0x02:
			return ELECTRICITY
		case 0x03:
			return GAS
		case 0x04:
			return HEAT_OUTLET
		case 0x05:
			return STREAM
		case 0x06:
			return HOT_WATER
		case 0x07:
			return WATER
		case 0x08:
			return HEAT_COST_ALLOC
		case 0x09:
			return COMPRESSED_AIR
		case 0x0A:
			return COLLING_RET_OUTLET
		case 0x0B:
			return COLLING_FLOW_INLET
		case 0x0C:
			return HEAT_INLET
		case 0x0D:
			return HEAT_COOLING
		case 0x0E:
			return BUS_SYSTEM
		case 0x0F:
			return UNKNOWN_MEDIUM
		case 0x16:
			return COLD_WATER
		case 0x17:
			return DUAL_WATER
		case 0x18:
			return PRESSURE
		case 0x19:
			return A_D_CONVERTER
		}

		return RESERVED

	}(answer[14])

	m.accessnumber = int(answer[15])
	m.status = int(answer[16])
	m.signature = int(answer[16])

}

func (m *Mbus) ParseDataRecords(answer []byte) {

	if len(answer) < 4 {
		return
	}

	dr := &DataRecord{}

	point := 0

	dif := answer[0]
	point += count_difs(answer[point:], 0)

	vif := answer[point]
	len_vib := count_vifs(answer[point:], 0)

	var vife byte = 0x00
	if len_vib > 1 {
		vife = answer[point+1]
	}

	point += len_vib

	parameter_, conversion, unit := parse_vif(vif, vife)

	len_data := get_len_data(dif)

	if len(answer) < point+len_data {
		return
	}

	value_b := answer[point : point+len_data]
	point += len_data

	value := ""

	if parameter_ == PARAMETER_DATE || parameter_ == PARAMETER_DATETIME {
		value = parseDate(value_b)
	} else {

		if int(dif&0xF) >= 8 {
			value = fmt.Sprintf("%d", convertBCD(value_b))
		} else {

			value = fmt.Sprintf("%d", convert_int(value_b))
			/*

				TODO доделать

					if n_ != 0 {
						format := "%" + fmt.Sprintf(".%df", n_)
						value = fmt.Sprintf(format, float64(convert_int(value_b)))
					} else {
						value = fmt.Sprintf("%d", convert_int(value_b))
					}
			*/
		}
	}

	dr.parameter = parameter_
	dr.value = value
	dr.unit = unit
	dr.conversion = conversion

	m.data = append(m.data, dr)

	m.ParseDataRecords(answer[point:])

}

// длина DIB
func count_difs(answer []byte, count int) int {
	return count_extends(answer, count)
}

// длина VIB
func count_vifs(answer []byte, count int) int {
	return count_extends(answer, count)
}

// Сколько раширяющих байт заниет заголовок?
func count_extends(answer []byte, count int) int {
	if count == 10 {
		return count
	}
	if len(answer) < 1 {
		return 0
	}
	if count == 0 {
		count += 1
	}
	if answer[0]>>7 == 1 {
		count += 1
		count_extends(answer[1:], count)
	}
	return count
}

func convertBCD(data []byte) int {
	str := ""
	for i := len(data) - 1; i >= 0; i-- {
		str += fmt.Sprintf("%02X", data[i])
	}
	integer, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return integer
}

func convert_int(int_data []byte) int {

	if len(int_data) == 0 {
		return 0
	}

	value := 0

	neg := int_data[len(int_data)-1] & 0x80

	for i := len(int_data); i > 0; i-- {
		if neg == 1 {
			value = (value << 8) + (int(int_data[i-1]) ^ 0xFF)
		} else {
			value = (value << 8) + int(int_data[i-1])
		}
	}
	if neg == 1 {
		value = value*-1 - 1
	}
	return value
}

func parseDate(t_data []byte) string {

	var sec, min, hour, day, mon, year int

	if len(t_data) == 6 {
		if t_data[1]&0x80 == 0 {
			sec = int(t_data[0]) & 0x3F
			min = int(t_data[1]) & 0x3F
			hour = int(t_data[2]) & 0x1F
			day = int(t_data[3]) & 0x1F
			mon = (int(t_data[4]) & 0x0F)
			year = 2000 + (((int(t_data[3]) & 0xE0) >> 5) | ((int(t_data[4]) & 0xF0) >> 1))

			return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, mon, day, hour, min, sec)

		}
	} else if len(t_data) == 4 {
		if t_data[0]&0x80 == 0 {
			min = int(t_data[0]) & 0x3F
			hour = int(t_data[1]) & 0x1F
			day = int(t_data[2]) & 0x1F
			mon = (int(t_data[3]) & 0x0F)
			year = 2000 + (((int(t_data[2]) & 0xE0) >> 5) | ((int(t_data[3]) & 0xF0) >> 1))
			return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, mon, day, hour, min)
		}
	} else if len(t_data) == 2 {
		day = int(t_data[0]) & 0x1F
		mon = (int(t_data[1]) & 0x0F)
		year = 2000 + (((int(t_data[0]) & 0xE0) >> 5) | ((int(t_data[1]) & 0xF0) >> 1))
		return fmt.Sprintf("%d-%02d-%02d", year, mon, day)
	}

	return ""
}

/*
	Разбираем VIF код
	@return parameter
	@return conversion
	@return unit
*/
func parse_vif(vif byte, vife byte) (string, string, string) {

	if vif == 0xFD || vif == 0xFB {

		if vife == 0x08 || vife == 0x88 {
			return PARAMETER_ACCESS_NUMBER, "", ""
		} else if vife == 0x09 || vife == 0x89 {
			return PARAMETER_MEDIUM, "", ""
		} else if vife == 0x0A || vife == 0x8A {
			return PARAMETER_MANUFACTURER, "", ""
		} else if vife == 0x0B || vife == 0x8B {
			return PARAMETER_IS_IDENTIFICATION, "", ""
		} else if vife == 0x0C || vife == 0x8C {
			return PARAMETER_MODEL_VERSION, "", ""
		} else if vife == 0x0D || vife == 0x8D {
			return PARAMETER_HARDWARE_VERSION, "", ""
		} else if vife == 0x0E || vife == 0x8E {
			return PARAMETER_FINDWARE_VERSION, "", ""
		} else if vife == 0x0F || vife == 0x8E {
			return PARAMETER_SOFTWARE_VERSION, "", ""
		} else if vife == 0x16 {
			return PARAMETER_PASSWORD, "", ""
		} else if vife == 0x17 || vife == 0x97 {
			return PARAMETER_ERROR_FLAG, "", ""
		} else if vife == 0x10 {
			return PARAMETER_CUSTOMER_LOCATION, "", ""
		} else if vife == 0x11 {
			return PARAMETER_CUSTOMER, "", ""
		} else if vife == 0x1A {
			return PARAMETER_DIGITAL_OUTPUT, "", ""
		} else if vife == 0x1B {
			return PARAMETER_DIGITAL_INPUT, "", ""
		} else if vife == 0x40 {
			return PARAMETER_V, "", ""
		} else if vife == 0x50 {
			return PARAMETER_A, "", ""
		} else {
			return PARAMETER_UNRECOGNIZED, "", ""
		}
		return PARAMETER_CUSTOM_VIF, "", ""
	}

	switch vif & 0x7F /* ignore the extension bit in this selection */ {

	// E000 0nnn Energy 10(nnn-3) W
	case 0x00, 0x00 + 1, 0x00 + 2, 0x00 + 3,
		0x00 + 4, 0x00 + 5, 0x00 + 6, 0x00 + 7:
		return PARAMETER_ENERGY, mbus_unit_prefix(int(vif & 0x07)), "W"

		// 0000 1nnn          Energy       10(nnn)J     (0.001kJ to 10000kJ)
	case 0x08, 0x08 + 1, 0x08 + 2, 0x08 + 3,
		0x08 + 4, 0x08 + 5, 0x08 + 6, 0x08 + 7:
		// mbus_unit_prefix(int(vif & 0x07))
		return PARAMETER_ENERGY, mbus_unit_prefix(int(vif & 0x07)), "J"

	// E001 1nnn Mass 10(nnn-3) kg 0.001kg to 10000kg
	case 0x18, 0x18 + 1, 0x18 + 2, 0x18 + 3,
		0x18 + 4, 0x18 + 5, 0x18 + 6, 0x18 + 7:
		return PARAMETER_MASS, mbus_unit_prefix(int(vif&0x07) - 3), "kg"

	// E010 1nnn Power 10(nnn-3) W 0.001W to 10000W
	case 0x28, 0x28 + 1, 0x28 + 2, 0x28 + 3,
		0x28 + 4, 0x28 + 5, 0x28 + 6, 0x28 + 7:
		return PARAMETER_POWER, mbus_unit_prefix(int(vif&0x07) - 3), "W"

	// E011 0nnn Power 10(nnn) J/h 0.001kJ/h to 10000kJ/h
	case 0x30, 0x30 + 1, 0x30 + 2, 0x30 + 3,
		0x30 + 4, 0x30 + 5, 0x30 + 6, 0x30 + 7:
		return PARAMETER_POWER, mbus_unit_prefix(int(vif & 0x07)), "J/h"

	// E001 0nnn Volume 10(nnn-6) m3 0.001l to 10000l
	case 0x10, 0x10 + 1, 0x10 + 2, 0x10 + 3,
		0x10 + 4, 0x10 + 5, 0x10 + 6, 0x10 + 7:
		return PARAMETER_VOLUME, mbus_unit_prefix(int(vif&0x07) - 6), "m^3"

	// E011 1nnn Volume Flow 10(nnn-6) m3/h 0.001l/h to 10000l/
	case 0x38, 0x38 + 1, 0x38 + 2, 0x38 + 3,
		0x38 + 4, 0x38 + 5, 0x38 + 6, 0x38 + 7:
		return PARAMETER_VOLUME_FLOW, mbus_unit_prefix(int(vif)&0x07 - 6), "m3/h"

	// E100 0nnn Volume Flow ext. 10(nnn-7) m3/min 0.0001l/min to 1000l/min
	case 0x40, 0x40 + 1, 0x40 + 2, 0x40 + 3,
		0x40 + 4, 0x40 + 5, 0x40 + 6, 0x40 + 7:
		return PARAMETER_VOLUME_FLOW, mbus_unit_prefix(int(vif&0x07) - 7), "m3/min"

	// E100 1nnn Volume Flow ext. 10(nnn-9) m3/s 0.001ml/s to 10000ml/
	case 0x48, 0x48 + 1, 0x48 + 2, 0x48 + 3,
		0x48 + 4, 0x48 + 5, 0x48 + 6, 0x48 + 7:
		return PARAMETER_VOLUME_FLOW, mbus_unit_prefix(int(vif&0x07) - 9), "m3/s"

	// E101 0nnn Mass flow 10(nnn-3) kg/h 0.001kg/h to 10000kg/
	case 0x50, 0x50 + 1, 0x50 + 2, 0x50 + 3,
		0x50 + 4, 0x50 + 5, 0x50 + 6, 0x50 + 7:
		return PARAMETER_MASS_FLOW, mbus_unit_prefix(int(vif&0x07) - 3), "kg/h"

	// E101 10nn Flow Temperature 10(nn-3) °C 0.001°C to 1°C
	case 0x58, 0x58 + 1, 0x58 + 2, 0x58 + 3:
		return PARAMETER_TEMP_FLOW, "dec", "C"

	// E101 11nn Return Temperature 10(nn-3) °C 0.001°C to 1°C
	case 0x5C, 0x5C + 1, 0x5C + 2, 0x5C + 3:
		return PARAMETER_TEMP_RETURN, "dec", "C"

	// E110 10nn Pressure 10(nn-3) bar 1mbar to 1000mbar
	case 0x68, 0x68 + 1, 0x68 + 2, 0x68 + 3:
		return PARAMETER_PRESSURE, mbus_unit_prefix(int(vif&0x03) - 3), "1mbar"

	// E010 00nn On Time
	// nn = 00 seconds
	// nn = 01 minutes
	// nn = 10   hours
	// nn = 11    days
	// E010 01nn Operating Time coded like OnTime
	// E111 00nn Averaging Duration coded like OnTime
	// E111 01nn Actuality Duration coded like OnTime
	case 0x20, 0x20 + 1, 0x20 + 2, 0x20 + 3,
		0x24, 0x24 + 1, 0x24 + 2, 0x24 + 3,
		0x70, 0x70 + 1, 0x70 + 2, 0x70 + 3,
		0x74, 0x74 + 1, 0x74 + 2, 0x74 + 3:

		// offset := 0

		if (vif & 0x7C) == 0x20 {
			return PARAMETER_ON_TIME, mbus_unit_prefix(0x00), ""
		} else if (vif & 0x7C) == 0x24 {
			return PARAMETER_OPERATING_TIME, mbus_unit_prefix(0x00), ""
		} else if (vif & 0x7C) == 0x70 {
			return PARAMETER_AVERAGING_DURATION, mbus_unit_prefix(0x00), ""
		} else {
			return PARAMETER_ACTUALITY_DURATION, mbus_unit_prefix(0x00), ""
		}

	// E110 110n Time Point
	// n = 0        date
	// n = 1 time & date
	// data type G
	// data type F
	case 0x6C,
		0x6C + 1:

		if vif&0x1 != 0 {
			return PARAMETER_DATETIME, "", ""
		}
		return PARAMETER_DATE, "", ""

	// E110 00nn    Temperature Difference   10(nn-3)K   (mK to  K)
	case 0x60, 0x60 + 1, 0x60 + 2, 0x60 + 3:
		return PARAMETER_TEMP_DIFF, mbus_unit_prefix(int(vif&0x03) - 3), "K"

	// E110 01nn External Temperature 10(nn-3) °C 0.001°C to 1°C
	case 0x64, 0x64 + 1, 0x64 + 2, 0x64 + 3:
		return PARAMETER_TEMP_EXTERNAL, mbus_unit_prefix(int(vif&0x03) - 3), "C"

	// E110 1110 Units for H.C.A. dimensionless
	case 0x6E:
		return PARAMETER_UNITS, "", ""

	// E110 1111 Reserved
	case 0x6F:
		return PARAMETER_RESERVED, "", ""

	// Custom VIF in the following string: never reached...
	case 0x7C:
		return PARAMETER_CUSTOM_VIF, "", ""

	// Fabrication No
	case 0x78:
		return PARAMETER_FABRICATION, "", ""

	// Bus Address
	case 0x7A:
		return PARAMETER_BUS_ADDR, "", ""

	// Manufacturer specific: 7Fh / FF
	case 0x7F,
		0xFF:

		return PARAMETER_MANUFACTURED_SPEC, "", ""

	default:
		return PARAMETER_UNDEFINED, "", ""
	}
}

func mbus_unit_prefix(exp int) string {
	switch exp {
	case -3:
		return "m"
	case -6:
		return "my"
	case 1:
		return "10 "
	case 2:
		return "100 "
	case 3:
		return "k"
	case 4:
		return "10 k"
	case 5:
		return "100 k"
	case 6:
		return "M"
	case 9:
		return "T"
	default:
		return fmt.Sprintf("1e%d ", exp)
	}
	return ""
}

func get_len_data(dif byte) int {
	switch dif & 0xF {
	case 0x1: // 1 byte integer (8 bit)
		return 1
	case 0x2: // 2 byte (16 bit)
		return 2
	case 0x3: // 3 byte integer (24 bit)
		return 3
	case 0x4: // 4 byte (32 bit)
		return 4
	case 0x5: // 4 Byte Real (32 bit)
		return 4
	case 0x6: // 6 byte (48 bit)
		return 5
	case 0x7: // 8 byte integer (64 bit)
		return 6
	case 0x9: // 2 digit BCD (8 bit)
		return 1
	case 0xA: // 4 digit BCD (16 bit)
		return 2
	case 0xB: // 6 digit BCD (24 bit)
		return 3
	case 0xC: // 8 digit BCD (32 bit)
		return 4
	case 0xE: // 12 digit BCD (48 bit)
		return 6
	case 0xF: // special functions
		return 6
	default:
		return 0
	}
}

func (m *Mbus) mbus_decode_manufacturer(a byte, b byte) string {
	var m_id int

	m_str := []byte{}

	m_id = m.mbus_data_int_decode([]byte{a, b}, 2)

	m_str = append(m_str, byte(((m_id>>10)&0x001F)+64))
	m_str = append(m_str, byte(((m_id>>5)&0x001F)+64))
	m_str = append(m_str, byte(((m_id)&0x001F)+64))
	m_str = append(m_str, 0x00)

	return string(m_str)
}

func (m *Mbus) mbus_data_int_decode(int_data []byte, ln int) int {
	value := 0
	if len(int_data) == 0 {
		return -1
	}

	neg := int_data[len(int_data)-1] & 0x80
	for i := ln; i > 0; i-- {
		if neg > 0 {
			value = (value << 8) + (int(int_data[i-1]) ^ 0xFF)
		} else {
			value = (value << 8) + int(int_data[i-1])
		}
	}
	if neg != 0 {
		value = (value * -1) - 1
	}

	return value
}
