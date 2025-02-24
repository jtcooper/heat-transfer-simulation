package main

import (
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ISystem interface {
	reset()
	step()
	commit(timeStep float64)
	getName() string
	getTemp() float64
	getData() map[string]*[]opts.LineData
}

type IFluidSystem interface {
	ISystem
	inputHeatCallback(heat float64)
}

type fluidSystem struct {
	name               string
	exposedSurfaceArea float64 // m^2; surface area exposed to the ambient environment
	ambientTemp        float64
	ambientHTC         float64
	fluidMass          float64
	temperature        float64 // internal fluid temp
	heatInComponents   []IComponent
	heatOutComponents  []IComponent
	// the step- prefix values need to be reset separately from the step function
	stepHeatIn  []float64
	stepHeatOut []float64
	powerData   map[string]*[]opts.LineData
}

func (fs fluidSystem) getName() string {
	return fs.name
}

func (fs fluidSystem) getTemp() float64 {
	return fs.temperature
}

func (fs fluidSystem) getData() map[string]*[]opts.LineData {
	return fs.powerData
}

func (fs *fluidSystem) addEnvironmentalConvectionHeatLossComponent() {
	fs.heatOutComponents = append(fs.heatOutComponents, &ambientConvectionHeatComponent{
		component: component{
			name: "Ambient Convection Heat Loss",
		},
		ambientHTC:  fs.ambientHTC,
		surfaceArea: fs.exposedSurfaceArea,
		currentTemp: func() float64 { return (*fs).temperature },
		ambientTemp: func() float64 { return (*fs).ambientTemp },
	})
}

func (fs *fluidSystem) addOutputHeatFluidComponent(output IFluidSystem, flowRate float64) {
	fs.heatOutComponents = append(fs.heatOutComponents,
		transferHeatComponentWrapper{
			component: component{
				name: "Heat Output",
			},
			wrappedComponent: heatCapacityFluidComponent{
				flowMass:     func() float64 { return flowRate },
				specificHeat: specificHeatWater,
				currentTemp:  func() float64 { return (*fs).temperature },
				outputTemp:   func() float64 { return output.getTemp() },
			},
			output: output,
		})
}

func (fs *fluidSystem) reset() {
	fs.stepHeatIn = []float64{}
	fs.stepHeatOut = []float64{}
}

func (fs *fluidSystem) step() {
	for _, comp := range fs.heatInComponents {
		if heatComp, ok := comp.(IHeatComponent); ok {
			q := heatComp.getHeat()
			fs.stepHeatIn = append(fs.stepHeatIn, q)
			fs.addDataPoint(heatComp.getName(), q)
		}
	}

	for _, comp := range fs.heatOutComponents {
		if heatComp, ok := comp.(IHeatComponent); ok {
			q := heatComp.getHeat()
			fs.stepHeatOut = append(fs.stepHeatOut, q)
			fs.addDataPoint(heatComp.getName(), q)
		}
		if fluidComp, ok := comp.(transferHeatComponentWrapper); ok {
			q := fluidComp.getHeat()
			fluidComp.transferHeat(q)
		}
	}
}

func (fs *fluidSystem) commit(timeStep float64) {
	// compute change in internal temperature
	// qᵢ - q₀ = qₛ
	heatStored := 0.0
	for _, q := range fs.stepHeatIn {
		heatStored += q
	}
	for _, q := range fs.stepHeatOut {
		heatStored -= q
	}

	// solve the heat capacity function for T₀
	// T₀ = q/ṁC + Tᵢ
	// note ṁ is in kg/s, so we need to factor in time passed
	fs.temperature = heatStored/((fs.fluidMass/timeStep)*specificHeatWater) + fs.temperature
}

func (fs *fluidSystem) inputHeatCallback(heat float64) {
	fs.stepHeatIn = append(fs.stepHeatIn, heat)
	fs.addDataPoint("Heat Input", heat)
}

func (fs *fluidSystem) addDataPoint(name string, value float64) {
	if fs.powerData == nil {
		fs.powerData = map[string]*[]opts.LineData{}
	}
	if _, ok := fs.powerData[name]; !ok {
		fs.powerData[name] = &[]opts.LineData{}
	}
	(*fs.powerData[name]) = append((*fs.powerData[name]), opts.LineData{Value: value})
}
