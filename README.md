# bird-watcher
Essentially it's a CCTV system, but it's meant specifically for watching birds at a feeder or in a nest box. It's designed to run on a Raspberry Pi 4, and as such does not feature GPU acceleration. However, it does make efficent use of all CPU cores available to it.

1 STRUCTURE
---------

The system is divided into several microservices, each of which is a seperate GoLand project. This structure allows one or more services to go offline without bringing down the entire system, and for different services to be optionally offloaded onto different physical machines for better performance. It also allows the system to be easily configured to incorporate multiple IP cameras.

All of the microservices are controlled with systemctl and they each have a config JSON file in /etc/bird-watcher-$VERSION.

  1.1 WEBCAM CONTROLLER SERVICE
  -----------------------------
  
  blah
  
  1.2 ACTIVITY DETECTOR SERVICE
  -----------------------------
  
  blah
  
  1.3 APPLICATION BACKEND SERVICE
  -------------------------------
  
  blah
  
  1.4 APPLICATION FRONTEND SERVICE
  --------------------------------
  
  blah
