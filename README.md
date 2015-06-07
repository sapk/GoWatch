# GoWatch [![Build Status](https://travis-ci.org/sapk/GoWatch.svg?branch=master)](https://travis-ci.org/sapk/GoWatch)
Simple network visualizer
##### Still under development
## [Goals](https://github.com/sapk/GoWatch/wiki/Detailled-goals)
- [ ] Analyze Network Unit
- [ ] SNMP collector
- [ ] Network Graphing
- [ ] Network Monitoring
- [ ] Network Map
- [ ] Network History/Research
- [ ] Service Monitoring
- [ ] Visual Configuration (simple)
- [ ] Config Backup & History
- [ ] Auto-Discovery
- [ ] Various Network Tools
- [ ] Organizations
- [ ] Location/Map
- [ ] Procedures (custom scripts)
- [ ] Equipement Web Console

For optimization there is two types of analysis :
  - Short ( very simple stats can be run in paralell)
    - Ping
    - SNMP Trap
    - ...
  - Long (Need more analysis so for less consumption need to be run sequentially (maybe 2-3 paralell can be changed in config))
    - SNMP
    - Log scan
    - Port scanning
    - Bandwith stats
    - ...

The admin dashboard show the medium interval between each scan of each type.

## Requirement
- librrd

## License

This project is under the MIT License. See the [LICENSE](https://github.com/sapk/GoWatch/blob/master/LICENSE) file for the full license text.
