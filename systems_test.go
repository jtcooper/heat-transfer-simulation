package main

import (
	"testing"
)

type mockComponent struct {
	name string
}

func (m mockComponent) getName() string {
	return m.name
}

type mockHeatComponent struct {
	mockComponent
	heat float64
}

func (m mockHeatComponent) getHeat() float64 {
	return m.heat
}

type mockTransferHeatComponentWrapper struct {
	mockComponent
	heat   float64
	output *mockFluidSystem
}

func (m mockTransferHeatComponentWrapper) getHeat() float64 {
	return m.heat
}

func (m *mockTransferHeatComponentWrapper) transferHeat(heat float64) {
	m.output.inputHeatCallback(heat)
}

func TestFluidSystem_Step(t *testing.T) {
	tests := []struct {
		name              string
		heatInComponents  []IComponent
		heatOutComponents []IComponent
		expectedHeatIn    []float64
		expectedHeatOut   []float64
	}{
		{
			name: "Single Heat In Component",
			heatInComponents: []IComponent{
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat Component"},
					heat:          1000.0,
				},
			},
			expectedHeatIn:  []float64{1000.0},
			expectedHeatOut: []float64{},
		},
		{
			name: "Single Heat Out Component",
			heatOutComponents: []IComponent{
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat Component"},
					heat:          500.0,
				},
			},
			expectedHeatIn:  []float64{},
			expectedHeatOut: []float64{500.0},
		},
		{
			name: "Multiple Heat In and Out Components",
			heatInComponents: []IComponent{
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat In Component 1"},
					heat:          1000.0,
				},
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat In Component 2"},
					heat:          2000.0,
				},
			},
			heatOutComponents: []IComponent{
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat Out Component 1"},
					heat:          500.0,
				},
				mockHeatComponent{
					mockComponent: mockComponent{name: "Heat Out Component 2"},
					heat:          1500.0,
				},
			},
			expectedHeatIn:  []float64{1000.0, 2000.0},
			expectedHeatOut: []float64{500.0, 1500.0},
		},
		{
			name: "Transfer Heat Component Wrapper",
			heatOutComponents: []IComponent{
				&mockTransferHeatComponentWrapper{
					mockComponent: mockComponent{name: "Transfer Heat Component"},
					heat:          300.0,
					output:        &mockFluidSystem{},
				},
			},
			expectedHeatIn:  []float64{},
			expectedHeatOut: []float64{300.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := fluidSystem{
				heatInComponents:  tt.heatInComponents,
				heatOutComponents: tt.heatOutComponents,
			}
			fs.step()
			if len(fs.stepHeatIn) != len(tt.expectedHeatIn) {
				t.Errorf("expected %v heat in components, got %v", len(tt.expectedHeatIn), len(fs.stepHeatIn))
			}
			for i, heat := range tt.expectedHeatIn {
				if fs.stepHeatIn[i] != heat {
					t.Errorf("expected heat in %v at index %v, got %v", heat, i, fs.stepHeatIn[i])
				}
			}
			if len(fs.stepHeatOut) != len(tt.expectedHeatOut) {
				t.Errorf("expected %v heat out components, got %v", len(tt.expectedHeatOut), len(fs.stepHeatOut))
			}
			for i, heat := range tt.expectedHeatOut {
				if fs.stepHeatOut[i] != heat {
					t.Errorf("expected heat out %v at index %v, got %v", heat, i, fs.stepHeatOut[i])
				}
			}
		})
	}
}

func TestFluidSystem_Commit(t *testing.T) {
	fs := fluidSystem{
		fluidMass:   2.0,
		temperature: 50.0,
		stepHeatIn:  []float64{1000.0},
		stepHeatOut: []float64{500.0},
	}
	fs.commit(1.0)
	expectedTemp := 50.0 + (500.0 / (2.0 * specificHeatWater))
	if fs.temperature != expectedTemp {
		t.Errorf("expected %v, got %v", expectedTemp, fs.temperature)
	}
}
