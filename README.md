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
