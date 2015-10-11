# Raspberry Pi Wireless Thermostat

Here's some code to make a Raspberry Pi control a Tower RFWRT wireless
thermostat. It's adaptable to other makes and models, but this is the one I've
got.

There's two bits. A C executable, called `boiler_control`, which when passed an
`on` or `off` argument, will send the appropriate signal to the heating control
unit.

There's a separate Go executable, called `thermostat`, that runs thermostat
control loop. It can be configured to read a file at regular intervals for the
target temperature, allowing the use of `cron` to schedule your heating.

`thermostat` can be cross compiled from any Go compiler, but it's probably
easiest to compile `boiler_control` on the Pi.
