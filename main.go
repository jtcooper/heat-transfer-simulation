package main

import (
	"fmt"
	"math"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// TODO: some values should be configurable
	outdoorAmbientTemp = 15.0  // Celsius
	indoorAmbientTemp  = 22.0  // Celsius
	outdoorHTC         = 15.0  // W/m^2*K
	indoorHTC          = 5.0   // W/m^2*K
	initialPanelTemp   = 30.0  // Celsius
	initialTankTemp    = 20.0  // Celsius
	massPanelFluid     = 10.0  // kg
	massTankFluid      = 250.0 // kg

	// constants
	solarIrradiance     = 1000.0  // W / m^2 ; varies with time of day
	specificHeatWater   = 4186.0  // J/(kg*K)
	flowRateCoefficient = 0.00286 // (kg/s) / m^2

	float64EqualityThreshold = 1e-9
)

var pumpFlowRatePerSecond = 0.0 // kg/s
var simulationRuntime = 3600

func main() {
	fmt.Println("Starting simulation...")
	panelSize := 2.0 // surface area: m^2
	// pumpFlowRatePerSecond = flowRateCoefficient * panelSize
	pumpFlowRatePerSecond = 0.5

	sp := solarPanel{
		fluidSystem: fluidSystem{
			name: "SolarPanel",
			// panel dimensions: 2m x 1m x 0.05m
			// if the panel is roof-mounted, all sides except the back are exposed to convection
			exposedSurfaceArea: 2.0*1.0 + 2*(2.0*0.05) + 2*(1.0*0.05),
			ambientTemp:        outdoorAmbientTemp,
			ambientHTC:         outdoorHTC,
			fluidMass:          massPanelFluid,
			temperature:        initialPanelTemp,
		},
		panelArea:       panelSize,
		panelEfficiency: 0.60, // typically ranges from 50-70%
	}

	st := storageTank{
		fluidSystem: fluidSystem{
			name: "StorageTank",
			// tank dimensions: 1.7m tall, 0.3m radius
			// A = 2πrh + πr^2 ; the side on the ground is insulated
			exposedSurfaceArea: 2*math.Pi*0.3*1.7 + math.Pi*math.Pow(0.3, 2),
			ambientTemp:        indoorAmbientTemp,
			ambientHTC:         indoorHTC,
			fluidMass:          massTankFluid,
			temperature:        initialTankTemp,
		},
	}

	// initialize the systems: hook up system outputs and inputs
	sp.initialize([]IFluidSystem{&st.fluidSystem})
	st.initialize([]IFluidSystem{&sp.fluidSystem})

	systems := []ISystem{&sp, &st}

	// setup plot arrays
	timeSeries := make([]float64, 0, simulationRuntime)
	systemsSeries := make(map[string][]opts.LineData)
	for _, sys := range systems {
		systemsSeries[sys.getName()] = make([]opts.LineData, 0, simulationRuntime)
	}

	// run the simulation
	timeStep := 1.0                         // seconds
	totalTime := float64(simulationRuntime) // 1 hour
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
	line.SetXAxis(timeSeries)
	for name, series := range systemsSeries {
		line.AddSeries(name, series)
	}
	plotLine(line, "temperatureSeries.html")

	// plot the results for each system's power values
	for _, sys := range systems {
		sysLine := charts.NewLine()
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
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	line.Render(f)
}
