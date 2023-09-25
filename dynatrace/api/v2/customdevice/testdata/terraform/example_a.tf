resource "dynatrace_custom_device" "#name#" {
  custom_device_id = "coffeeDeviceId"
  dns_names        = [ "coffee-machine.dynatrace.internal.com" ]
  group            = "myCustomDeviceGroup"
  props            = jsonencode({
      "coffee": "caffeinated"
  })  
  type         = "CUSTOM_DEVICE"
  display_name = "coffeeMachine"
  favicon_url  = "https://www.freefavicon.com/freefavicons/food/cup-of-coffee-152-78475.png"
  ip_addresses = [ "10.0.0.1" ]
  listen_ports = [80]
}