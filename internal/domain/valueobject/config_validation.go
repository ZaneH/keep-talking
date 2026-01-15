package valueobject

import (
	"errors"
	"fmt"
	"strings"
)

// Configuration limits
const (
	MinTimerSeconds      = 5
	MaxTimerSeconds      = 3600 // 1 hour
	MinStrikes           = 1
	MaxStrikes           = 10
	MinFaces             = 1
	MaxFaces             = 10
	MinRows              = 1
	MaxRows              = 4
	MinColumns           = 1
	MaxColumns           = 5
	MinModules           = 1
	MaxModulesTotal      = 60
	MaxModulesPerFaceCap = 15
	MinBatteriesAllowed  = 0
	MaxBatteriesAllowed  = 6
	MaxIndicatorsAllowed = 5
	MaxPortsAllowed      = 6
	MinLevel             = 1
	MaxLevel             = 10
)

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msg strings.Builder
	msg.WriteString("validation errors:")
	for _, err := range e {
		msg.WriteString("\n  - " + err.Error())
	}
	return msg.String()
}

func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

func ValidateBombConfig(config BombConfig) ValidationErrors {
	var errs ValidationErrors

	// Timer validation
	timerSeconds := int(config.Timer.Seconds())
	if timerSeconds < MinTimerSeconds {
		errs = append(errs, ValidationError{
			Field:   "timer",
			Message: fmt.Sprintf("must be at least %d seconds", MinTimerSeconds),
		})
	}
	if timerSeconds > MaxTimerSeconds {
		errs = append(errs, ValidationError{
			Field:   "timer",
			Message: fmt.Sprintf("cannot exceed %d seconds", MaxTimerSeconds),
		})
	}

	// Strikes validation
	if config.MaxStrikes < MinStrikes {
		errs = append(errs, ValidationError{
			Field:   "max_strikes",
			Message: fmt.Sprintf("must be at least %d", MinStrikes),
		})
	}
	if config.MaxStrikes > MaxStrikes {
		errs = append(errs, ValidationError{
			Field:   "max_strikes",
			Message: fmt.Sprintf("cannot exceed %d", MaxStrikes),
		})
	}

	// Faces validation
	if config.NumFaces < MinFaces {
		errs = append(errs, ValidationError{
			Field:   "num_faces",
			Message: fmt.Sprintf("must be at least %d", MinFaces),
		})
	}
	if config.NumFaces > MaxFaces {
		errs = append(errs, ValidationError{
			Field:   "num_faces",
			Message: fmt.Sprintf("cannot exceed %d", MaxFaces),
		})
	}

	// Grid validation
	if config.Rows < MinRows || config.Rows > MaxRows {
		errs = append(errs, ValidationError{
			Field:   "rows",
			Message: fmt.Sprintf("must be between %d and %d", MinRows, MaxRows),
		})
	}
	if config.Columns < MinColumns || config.Columns > MaxColumns {
		errs = append(errs, ValidationError{
			Field:   "columns",
			Message: fmt.Sprintf("must be between %d and %d", MinColumns, MaxColumns),
		})
	}

	// Module count validation
	if config.MinModules < MinModules {
		errs = append(errs, ValidationError{
			Field:   "min_modules",
			Message: fmt.Sprintf("must be at least %d", MinModules),
		})
	}

	totalPossibleSlots := config.NumFaces * config.Rows * config.Columns
	if config.MinModules > totalPossibleSlots {
		errs = append(errs, ValidationError{
			Field: "min_modules",
			Message: fmt.Sprintf("exceeds available slots (%d faces x %d rows x %d cols = %d)",
				config.NumFaces, config.Rows, config.Columns, totalPossibleSlots),
		})
	}
	if config.MinModules > MaxModulesTotal {
		errs = append(errs, ValidationError{
			Field:   "min_modules",
			Message: fmt.Sprintf("cannot exceed %d modules total", MaxModulesTotal),
		})
	}

	if config.MaxModulesPerFace > MaxModulesPerFaceCap {
		errs = append(errs, ValidationError{
			Field:   "max_modules_per_face",
			Message: fmt.Sprintf("cannot exceed %d", MaxModulesPerFaceCap),
		})
	}

	// Battery validation
	if config.MinBatteries < MinBatteriesAllowed || config.MaxBatteries > MaxBatteriesAllowed {
		errs = append(errs, ValidationError{
			Field:   "batteries",
			Message: fmt.Sprintf("must be between %d and %d", MinBatteriesAllowed, MaxBatteriesAllowed),
		})
	}
	if config.MinBatteries > config.MaxBatteries {
		errs = append(errs, ValidationError{
			Field:   "batteries",
			Message: "min_batteries cannot exceed max_batteries",
		})
	}

	// Indicator validation
	if config.MaxIndicatorCount > MaxIndicatorsAllowed {
		errs = append(errs, ValidationError{
			Field:   "max_indicator_count",
			Message: fmt.Sprintf("cannot exceed %d", MaxIndicatorsAllowed),
		})
	}

	// Port validation
	if config.PortCount > MaxPortsAllowed {
		errs = append(errs, ValidationError{
			Field:   "port_count",
			Message: fmt.Sprintf("cannot exceed %d", MaxPortsAllowed),
		})
	}

	return errs
}

func ValidateLevel(level int) error {
	if level < MinLevel || level > MaxLevel {
		return fmt.Errorf("level must be between %d and %d", MinLevel, MaxLevel)
	}
	return nil
}

func ValidateMission(mission Mission) error {
	if _, exists := MissionDefinitions[mission]; !exists {
		return errors.New("unknown mission")
	}
	return nil
}
