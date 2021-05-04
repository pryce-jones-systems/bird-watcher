# bird-watcher
Essentially it's a CCTV system, but it's meant specifically for watching birds at a feeder or in a nest box. It's designed to run on a Raspberry Pi 4, and as such does not feature GPU acceleration. However, it does make efficient use of all CPU cores available to it.

1 STRUCTURE
---------

The system is divided into several microservices, each of which is a seperate GoLand project. This structure allows one or more services to go offline without bringing down the entire system, and for different services to be optionally offloaded onto different physical machines for better performance. It also allows the system to be easily configured to incorporate multiple IP cameras.

All of the microservices are controlled with systemctl and they each have a config JSON file in /etc/bird-watcher-$VERSION.

  1.1 WEBCAM CONTROLLER SERVICE (WCS)
  -----------------------------------
  
  This one's pretty simple (it's only about 400 lines!) This microservice's job is to read raw data from a video device, convert it to JPEG, and stream it over HTTP. If you already have an IP camera, then you don't need to use this microservice because all it does is turn a USB webcam into an IP camera.
  
  The BlackJack library (github.com/blackjack/webcam) is used to interact with the Video4Linux susbsystem. This microservice is super-lightweight, so it can run on even a really underpowered device.
  
  1.2 ACTIVITY DETECTOR SERVICE (ADS)
  -----------------------------------
  
  This microservice's job is to monitor an MJPEG image stream over HTTP (the stream can be local or remote) and decide whether any activity that is worth recording is happening. If such activity is taking place, a video is recorded and placed in a directory where it is accessible to the ABS.
  
  There are a number of ways to detect activity - the simplest of course is just to detect movement - but the best way (at least for bird watching) is to use a machine learning technique to determine what is going on. This microservice has several methods of detecting activity (including machine learning). The method used and relevant thresholds are set in the config file.
  
  1.3 APPLICATION BACKEND SERVICE (ABS)
  -------------------------------------
  
  blah
  
  1.4 APPLICATION FRONTEND SERVICE (AFS)
  --------------------------------------
  
  blah
  
2 INSTALLATION
--------------

To build and install the project, run install.sh with root privileges.

Config files are located at /etc/bird-watcher$VERSION
Binaries are located at /bin/bird-watcher$VERSION

Any of the services can be executed from any directory, but placing the binaries in /bin gives every user access. If you want to play around with different configurations, launch a service from any directory and make sure there is a file called config.json in that directory. This config file will override the one in /etc/bird-watcher-$VERSION
