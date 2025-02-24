package main

import (
	"testing"

	"github.com/go-echarts/go-echarts/v2/opts"
)

func mockVariableIntegrator(value float64) variableIntegrator {
	return func() float64 {
		return value
	}
}

func TestAmbientConvectionHeatComponent(t *testing.T) {
	component := ambientConvectionHeatComponent{
		component:   component{name: "Ambient Convection"},
		ambientHTC:  10.0,
		surfaceArea: 5.0,
		currentTemp: mockVariableIntegrator(100.0),
		ambientTemp: mockVariableIntegrator(20.0),
	}

	expectedHeat := 4000.0 // q = hAΔT = 10 * 5 * (100 - 20)
	if heat := component.getHeat(); heat != expectedHeat {
		t.Errorf("expected %v, got %v", expectedHeat, heat)
	}
}

func TestHeatAborptionComponent(t *testing.T) {
	component := heatAborptionComponent{
		component:         component{name: "Heat Absorption"},
		efficiency:        0.8,
		incidentRadiation: mockVariableIntegrator(1000.0),
		surfaceArea:       5.0,
	}

	expectedHeat := 4000.0 // q = ηIA = 0.8 * 1000 * 5
	if heat := component.getHeat(); heat != expectedHeat {
		t.Errorf("expected %v, got %v", expectedHeat, heat)
	}
}

func TestHeatCapacityFluidComponent(t *testing.T) {
	component := heatCapacityFluidComponent{
		component:    component{name: "Heat Capacity Fluid"},
		flowMass:     mockVariableIntegrator(2.0),
		specificHeat: 4.18,
		currentTemp:  mockVariableIntegrator(80.0),
		outputTemp:   mockVariableIntegrator(60.0),
	}

	expectedHeat := 167.2 // q = ṁCΔT = 2 * 4.18 * (80 - 60)
	if heat := component.getHeat(); heat != expectedHeat {
		t.Errorf("expected %v, got %v", expectedHeat, heat)
	}
}

type mockFluidSystem struct {
	receivedHeat float64
}

func (mockFluidSystem) reset()                               {}
func (mockFluidSystem) step()                                {}
func (mockFluidSystem) commit(timeStep float64)              {}
func (mockFluidSystem) getName() string                      { return "Mock Fluid System" }
func (mockFluidSystem) getTemp() float64                     { return 0.0 }
func (mockFluidSystem) getData() map[string]*[]opts.LineData { return map[string]*[]opts.LineData{} }
func (m *mockFluidSystem) inputHeatCallback(heat float64) {
	m.receivedHeat = heat
}

func TestTransferHeatComponentWrapper(t *testing.T) {
	mockOutput := &mockFluidSystem{}
	wrappedComponent := heatAborptionComponent{
		component:         component{name: "Wrapped Heat Absorption"},
		efficiency:        0.8,
		incidentRadiation: mockVariableIntegrator(1000.0),
		surfaceArea:       5.0,
	}

	component := transferHeatComponentWrapper{
		component:        component{name: "Transfer Heat Wrapper"},
		wrappedComponent: wrappedComponent,
		output:           mockOutput,
	}

	expectedHeat := 4000.0 // q = ηIA = 0.8 * 1000 * 5
	if heat := component.getHeat(); heat != expectedHeat {
		t.Errorf("expected %v, got %v", expectedHeat, heat)
	}

	component.transferHeat(expectedHeat)
	if mockOutput.receivedHeat != expectedHeat {
		t.Errorf("expected %v, got %v", expectedHeat, mockOutput.receivedHeat)
	}
}
