package main

import (
	"errors"
	"os"
	"strconv"
)

// default values can be overriden with environment variables
const (
	outdoorAmbientTemp = 15.0   // Celsius
	indoorAmbientTemp  = 22.0   // Celsius
	outdoorHTC         = 15.0   // W/m^2*K
	indoorHTC          = 5.0    // W/m^2*K
	panelTemp          = 30.0   // Celsius
	tankTemp           = 20.0   // Celsius
	panelFluidMass     = 10.0   // kg
	tankFluidMass      = 250.0  // kg
	solarIrradiance    = 1000.0 // W/m^2
	pumpFlowRate       = 0.2    // kg/s
	panelSize          = 2.0    // m^2
	panelEfficiency    = 0.6
	durationHours      = 1.0 // hr
)

type config struct {
	outdoorAmbientTemp float64
	indoorAmbientTemp  float64
	outdoorHTC         float64
	indoorHTC          float64
	panelTemp          float64
	tankTemp           float64
	panelFluidMass     float64
	tankFluidMass      float64
	solarIrradiance    float64
	pumpFlowRate       float64
	panelSize          float64
	panelEfficiency    float64
	durationHours      float64
}

func initializeConfig() config {
	config := config{
		outdoorAmbientTemp: outdoorAmbientTemp,
		indoorAmbientTemp:  indoorAmbientTemp,
		outdoorHTC:         outdoorHTC,
		indoorHTC:          indoorHTC,
		panelTemp:          panelTemp,
		tankTemp:           tankTemp,
		panelFluidMass:     panelFluidMass,
		tankFluidMass:      tankFluidMass,
		solarIrradiance:    solarIrradiance,
		pumpFlowRate:       pumpFlowRate,
		panelSize:          panelSize,
		panelEfficiency:    panelEfficiency,
		durationHours:      durationHours,
	}

	var err error
	if val := os.Getenv("OUTDOOR_TEMP"); val != "" {
		config.outdoorAmbientTemp, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("INDOOR_TEMP"); val != "" {
		config.indoorAmbientTemp, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("OUTDOOR_HTC"); val != "" {
		config.outdoorHTC, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("INDOOR_HTC"); val != "" {
		config.indoorAmbientTemp, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("PANEL_TEMP"); val != "" {
		config.panelTemp, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("TANK_TEMP"); val != "" {
		config.tankTemp, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("PANEL_WATER_MASS"); val != "" {
		config.panelFluidMass, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("TANK_WATER_MASS"); val != "" {
		config.tankFluidMass, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("SOLAR_IRRADIANCE"); val != "" {
		config.solarIrradiance, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("PUMP_FLOW_RATE"); val != "" {
		config.pumpFlowRate, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("PANEL_SIZE"); val != "" {
		config.panelSize, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("PANEL_EFFICIENCY"); val != "" {
		config.panelEfficiency, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	if val := os.Getenv("DURATION_HOURS"); val != "" {
		config.durationHours, err = strconv.ParseFloat(val, 64)
		handleParseEnvError(err)
	}
	return config
}

func handleParseEnvError(err error) {
	if err != nil {
		panic(errors.New("could not parse environment variable to float"))
	}
}
