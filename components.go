// a component is some physics component that produces one or more values
// components have a combination of constant and variable values
// components are as dumb as possible: they only compute values, and if needed, they have an output system for transferring power.
// To make them both generic and allow them to depend on variable values, components define variableIntegrator methods which the caller defines
package main

type variableIntegrator func() float64

// component can be used for all types of components that matter to a system.
// For this simple simulation, the only two types are heat and fluid components.
type IComponent interface {
	getName() string
}

type component struct {
	name string
}

func (c component) getName() string {
	return c.name
}

// IHeatComponent defines a method for getting heat rate
type IHeatComponent interface {
	IComponent
	getHeat() (heat float64)
}

type ambientConvectionHeatComponent struct {
	component
	ambientHTC  float64
	surfaceArea float64
	currentTemp variableIntegrator
	ambientTemp variableIntegrator
}

func (c ambientConvectionHeatComponent) getHeat() float64 {
	// q = hAΔT
	return c.ambientHTC * c.surfaceArea * (c.currentTemp() - c.ambientTemp())
}

type heatAborptionComponent struct {
	component
	efficiency        float64
	incidentRadiation variableIntegrator
	surfaceArea       float64
}

func (c heatAborptionComponent) getHeat() float64 {
	// q = ηIA
	return c.efficiency * c.incidentRadiation() * c.surfaceArea
}

type heatCapacityFluidComponent struct {
	component
	flowMass     variableIntegrator // kg
	specificHeat float64
	currentTemp  variableIntegrator
	outputTemp   variableIntegrator
}

func (c heatCapacityFluidComponent) getHeat() float64 {
	// q = ṁCΔT
	return c.flowMass() * c.specificHeat * (c.currentTemp() - c.outputTemp())
}

// transferHeatComponentWrapper is more complex than a simple component.
// It wraps an IHeatComponent and can be used to transfer heat to another system
type transferHeatComponentWrapper struct {
	component
	wrappedComponent IComponent
	output           IFluidSystem
}

func (oc transferHeatComponentWrapper) getHeat() float64 {
	if hc, ok := oc.wrappedComponent.(IHeatComponent); ok {
		return hc.getHeat()
	} else {
		panic("getHeat called for non-heatComponent")
	}
}

func (oc *transferHeatComponentWrapper) transferHeat(heat float64) {
	oc.output.inputHeatCallback(heat)
}
