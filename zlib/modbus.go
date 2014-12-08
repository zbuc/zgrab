package zlib

import "encoding/json"

type ModbusEvent struct {
	Response []byte
}

var ModbusEventType = EventType{
	TypeName:         CONNECTION_EVENT_MODBUS,
	GetEmptyInstance: func() EventData { return new(ModbusEvent) },
}

func (m *ModbusEvent) GetType() EventType {
	return ModbusEventType
}

type encodedModbusEvent struct {
	Response []byte `json:"response"`
}

func (m *ModbusEvent) MarshalJSON() ([]byte, error) {
	e := encodedModbusEvent{
		Response: m.Response,
	}
	return json.Marshal(&e)
}

func (m *ModbusEvent) UnmarshalJSON(b []byte) error {
	e := new(encodedModbusEvent)
	if err := json.Unmarshal(b, e); err != nil {
		return err
	}
	m.Response = e.Response
	return nil
}

type FunctionCode byte
type ExceptionFunctionCode byte
type ExceptionCode byte

type ModbusRequest struct {
	Function FunctionCode
	Data     []byte
}

type ModbusResponse struct {
	Function FunctionCode
	Data     []byte
}

type ModbusException struct {
	Function      ExceptionFunctionCode
	ExceptionType byte
}

func (e ExceptionFunctionCode) FunctionCode() FunctionCode {
	code := byte(e) & byte(0x79)
	return FunctionCode(code)
}

func (c FunctionCode) ExceptionFunctionCode() ExceptionFunctionCode {
	code := byte(c) | byte(0x80)
	return ExceptionFunctionCode(code)
}

var ModbusReadDeviceIDRequest = []byte{0x2B, 0x0E, 0x01, 0x00}
