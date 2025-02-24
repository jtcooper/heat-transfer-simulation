package main

import (
	"fmt"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	float64EqualityThreshold = 1e-9
	specificHeatWater        = 4186.0 // J/(kg*K)
)

func main() {
	config := initializeConfig()

	sp := solarPanel{
		fluidSystem: fluidSystem{
			name: "SolarPanel",
			// for simplicity, we'll ignore panel depth and treat the panel as if it is mounted on the roof
			exposedSurfaceArea: config.panelSize,
			ambientTemp:        config.outdoorAmbientTemp,
			ambientHTC:         config.outdoorHTC,
			fluidMass:          config.panelFluidMass,
			temperature:        config.panelTemp,
		},
		panelArea:       config.panelSize,
		panelEfficiency: config.panelEfficiency,
		solarIrradiance: config.solarIrradiance,
	}

	st := storageTank{
		fluidSystem: fluidSystem{
			name: "StorageTank",
			// for simplicty, tank dimensions aren't configurable
			// tank dimensions: 1.7m tall, 0.3m radius
			// A = 2πrh + πr^2 , where the side touching the ground is insulated
			exposedSurfaceArea: 2*math.Pi*0.3*1.7 + math.Pi*math.Pow(0.3, 2),
			ambientTemp:        config.indoorAmbientTemp,
			ambientHTC:         config.indoorHTC,
			fluidMass:          config.tankFluidMass,
			temperature:        config.tankTemp,
		},
	}

	// initialize the systems: hook up system outputs and inputs
	sp.initialize([]IFluidSystem{&st.fluidSystem}, config.pumpFlowRate)
	st.initialize([]IFluidSystem{&sp.fluidSystem}, config.pumpFlowRate)

	systems := []ISystem{&sp, &st}
	totalTime := config.durationHours * 60 * 60
	simulationRuntime := int(totalTime)

	// setup plot arrays
	timeSeries := make([]float64, 0, simulationRuntime)
	systemsSeries := make(map[string][]opts.LineData)
	for _, sys := range systems {
		systemsSeries[sys.getName()] = make([]opts.LineData, 0, simulationRuntime)
	}

	// run the simulation
	fmt.Println("Starting simulation...")
	timeStep := 1.0 // seconds
	for t := 0.0; t < totalTime+float64EqualityThreshold; t += timeStep {
		timeSeries = append(timeSeries, t)
		for _, sys := range systems {
			sys.reset()
		}
		for _, sys := range systems {
			sys.step()
		}
		for _, sys := range systems {
			sys.commit(timeStep)
			systemsSeries[sys.getName()] = append(systemsSeries[sys.getName()], opts.LineData{Value: sys.getTemp()})
		}
	}

	// plot the temperature results
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Temperature",
		Subtitle: "Solar Panel temperature vs. Storage Tank temperature over time",
	}))
	line.SetXAxis(timeSeries)
	for name, series := range systemsSeries {
		line.AddSeries(name, series)
	}
	plotLine(line, "TemperatureSeries.html")

	// plot the results for each system's power values
	for _, sys := range systems {
		sysLine := charts.NewLine()
		sysLine.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
			Title: sys.getName(),
		}))
		sysLine.SetXAxis(timeSeries)
		for name, sysSeries := range sys.getData() {
			sysLine.AddSeries(name, *sysSeries)
		}
		plotLine(sysLine, sys.getName()+"Series.html")
	}

	fmt.Println("Complete.")
}

func plotLine(line *charts.Line, fileName string) {
	// render to an HTML file
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	line.Render(f)
}
