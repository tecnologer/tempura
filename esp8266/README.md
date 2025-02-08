# ESP8266: Sensors Configuration

## SW-AGUA-INOX

A steel sensor of the water level with a voltage range of 0 – 250.

SW-AGUA-INOX | Pin	Connection
----|----
One Wire |	D1 (GPIO5) or any digital pin
Other Wire |	GND

### Buy at: 
 - [Mercado Libre][1]

## SHT30

Soil Temperature and Humidity Sensor High Precision Probe with I2C Output with 1.5m Cable DC3.3V

SHT30 Pin |	ESP8266 Pin |	Description
----|----|----
VCC	| 3.3V |	Power supply
GND	| GND	| Ground
SDA	| D2 (GPIO4)	I2C | Data Line
SCK (SCL)	| D1 (GPIO5)	I2C|  Clock Line

### Buy at:
 - [Mercado Libre][2]

## Battery LiPo 1 Cell

Lipo battery 3.7v to 1200mah.

Component |	Connection
---|----
Battery Positive |	Connect to R1 (30kΩ)
R1 (30kΩ)	| Connect to R2 (10kΩ)
Between R1 and R2 | 	Connect to A0 pin on ESP8266
R2 (10kΩ)	| Connect to GND
Battery Ground	| Connect to GND on ESP8266

```
  [ Battery + ] ----[ R1 (33kΩ or 39kΩ) ]----+----[ A0 on ESP8266 ]
                                             |
                                         [ R2 (10kΩ) ]
                                             |
  [ Battery - ] ----------------------------[ GND on ESP8266 ]
```

### Buy at:
  - [Mercado Libre][3]

[1]: https://articulo.mercadolibre.com.mx/MLM-593332982-switch-sensor-nivel-agua-metal-arduino-pic-avr-raspberry-_JM
[2]: https://www.mercadolibre.com.mx/sensor-de-humedad-del-suelo-sonda-de-temperatura-alta-prec/p/MLM2000952700
[3]: https://articulo.mercadolibre.com.mx/MLM-3264635028-bateria-lipo-37v-1200mah-3cables-_JM
