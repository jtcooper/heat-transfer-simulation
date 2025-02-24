# Heat Transfer Simulation

The repo includes an executable, so installing Go is not required to run the simulation.

Steps to run:

```
git clone https://github.com/jtcooper/heat-transfer-simulation.git
cd heat-transfer-simulation
./heat-transfer-simulation
```

The simulation outputs three HTML files:
* `TemperatureSeries.html` plots temperature of the solar panel and storage tank
* `SolarPanelSeries.html` plots the heat transfer values for the solar panel
* `StorageTankSeries.html` plots the heat transfer values for the storage tank

Select a series' name in the legend to toggle visibility.

## Adjusting simulation parameters

Most of the simulation's parameters can be adjusted through environment variables. The following shows all configurable variables with their default values:

```
OUTDOOR_TEMP=15 \
INDOOR_TEMP=22 \
OUTDOOR_HTC=15 \
INDOOR_HTC=5 \
PANEL_TEMP=30 \
TANK_TEMP=20 \
PANEL_WATER_MASS=10 \
TANK_WATER_MASS=250 \
SOLAR_IRRADIANCE=1000 \
PUMP_FLOW_RATE=0.2 \
PANEL_SIZE=2 \
PANEL_EFFICIENCY=0.6 \
DURATION_HOURS=1 \
./heat-transfer-simulation
```

## Design considerations

This simple solution involves two systems:
* Solar panel
* Storage tank

I chose to ignore the pipes, so heat transfer occurs directly between the solar panel and storage tank (considerations: pipe length, varying insulation through walls or exposed pipes)

There are a lot of heat components that could have been considered in this simulation. For simplicity, I just chose a few:
* Solar panel:
  * Incident solar radiation
  * Convection heat loss to the ambient environment
  * Heat transfer to and from the storage tank
* Storage tank:
  * Convection heat loss to the ambient environment
  * Heat transfer to and from the solar panel

I assumed the panel is in an outdoor environment and the storage tank is in an indoor environment. This means only the solar panel is exposed to solar irradiance, and they use different parameters for convection.

Considerations:
* I chose to ignore conduction heat loss through walls, touching objects, etc.
* Each system has uniform temperature, ignoring components such as thermal stratification in the storage tank.
* I chose to ignore heat loss due to radiation. The temperature differences are fairly small so this would not have been impactful.
* Solar irradiance typically varies throughout the day, but I chose to keep it constant for simplicity.
* The water flow from the pump is set to a constant rate that's applied to the entire system.
* The focus of this exercise is heat transfer, so I ignored other components such as pressure.
