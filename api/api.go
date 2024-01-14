package api

import (
	"errors"
	"net"
	"time"
)

const (
	API_Key_Command  = "_command"
	API_Key_Address  = "_address"
	API_Key_Duration = "_duration"

	API_Command_StartEmitter = "start_emitter"
	API_Command_StopEmitter  = "stop_emitter"

	Key_Timestamp    = "timestamp"
	Key_Hostname     = "hostname"
	Key_DirUsrHome   = "dir_usr_home"
	Key_DirUsrConfig = "dir_usr_config"
	Key_DirUsrCache  = "dir_usr_cache"
)

var (
	ErrCommandKeyNotFound  = errors.New("'_command' key not found")
	ErrCommandKeyNotString = errors.New("'_command' is not a string")
	ErrCommandUnknown      = errors.New("unknown '_command'")

	ErrAddressKeyNotFound    = errors.New("no '_address' key found")
	ErrAddressKeyNotString   = errors.New("'_address' is not a string")
	ErrAddressKeyNotResolved = errors.New("cannot resolve '_address'")

	ErrDurationKeyNotFound  = errors.New("'_duration' key not found")
	ErrDurationKeyNotString = errors.New("'_duration' is not a string")
	ErrDurationKeyNotParsed = errors.New("cannot parse '_duration'")

	ErrInvalidAddressKey = errors.New("invalid '_address'")
)

func IsValid(form map[string]any) error {

	commandVal, ok := form[API_Key_Command]
	if !ok {
		return ErrCommandKeyNotFound
	}

	commandStr, ok := commandVal.(string)
	if !ok {
		return ErrCommandKeyNotString
	}

	switch commandStr {
	case API_Command_StartEmitter:
		return isValidStartEmitter(form)
	case API_Command_StopEmitter:
		return isValidStopEmitter(form)
	default:
		return ErrCommandUnknown
	}
}

func isValidStartEmitter(form map[string]any) error {
	if err := isValidStopEmitter(form); err != nil {
		return err
	}

	durationVal, ok := form[API_Key_Duration]
	if !ok {
		return ErrDurationKeyNotFound
	}

	durationStr, ok := durationVal.(string)
	if !ok {
		return ErrDurationKeyNotString
	}

	if _, err := time.ParseDuration(durationStr); err != nil {
		return ErrDurationKeyNotParsed
	}

	return nil
}

func isValidStopEmitter(form map[string]any) error {
	addressVal, ok := form[API_Key_Address]
	if !ok {
		return ErrAddressKeyNotFound
	}

	addressStr, ok := addressVal.(string)
	if !ok {
		return ErrAddressKeyNotString
	}

	if _, err := net.ResolveUDPAddr("udp", addressStr); err != nil {
		return ErrAddressKeyNotResolved
	}

	return nil
}
