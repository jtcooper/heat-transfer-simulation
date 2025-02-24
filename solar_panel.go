// solar panel: with serpentine design: increases efficiency, but requires a low flow rate
package main

type solarPanel struct {
	fluidSystem
	panelArea       float64
	panelEfficiency float64
	solarIrradiance float64
}

func (sp *solarPanel) initialize(fluidOutputs []IFluidSystem, flowRate float64) {
	// include all the power components involved in this system
	sp.heatInComponents = []IComponent{
		heatAborptionComponent{
			component: component{
				name: "Incident Radiation",
			},
			efficiency:        sp.panelEfficiency,
			incidentRadiation: func() float64 { return (*sp).solarIrradiance },
			surfaceArea:       sp.panelArea,
		},
	}

	sp.heatOutComponents = []IComponent{}
	sp.addEnvironmentalConvectionHeatLossComponent()
	for _, output := range fluidOutputs {
		sp.addOutputHeatFluidComponent(output, flowRate)
	}
}
