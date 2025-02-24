package main

type storageTank struct {
	fluidSystem
}

func (st *storageTank) initialize(fluidOutputs []IFluidSystem, flowRate float64) {
	// include all the power components involved in this system
	st.heatOutComponents = []IComponent{}
	st.addEnvironmentalConvectionHeatLossComponent()
	for _, output := range fluidOutputs {
		st.addOutputHeatFluidComponent(output, flowRate)
	}
}
