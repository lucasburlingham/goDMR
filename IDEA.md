Complete rewrite of pi-star and wpsd. Use go to build a single binary that does everything and eliminates the need for multiple services, while being lightweight and fast. 

Target hardware: Original Raspberry Pi Zero with 512MB RAM.

Target coverage of digital voice modes: DMR. I don't have the hardware to test other modes, lets just get DMR working first and if anyone wants to help test other modes we can add them later.

Target features:
- Web interface for configuration and monitoring.
- Support for multiple DMR networks (BrandMeister, DMR+, etc).
- Hotspot support with low latency.
- Support for multiple radios (USB, GPIO, etc).
- Automatic updates and easy installation process.
- Logging and diagnostics tools.

No: 
- Endless feature bloat. Keep it simple and focused on core functionality.
- Endless configuration options. Sensible defaults and minimal setup.
- Multi-user support. Single user only for simplicity.

This is:
- An appliance. Plug it in, configure it via web interface, use it.
- GPLv3 licensed. Open source and free to use.