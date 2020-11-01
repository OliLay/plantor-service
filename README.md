# plantor-service - backend for [plantor](https://github.com/OliLay/plantor)

Provides basic services like
- MQTT broker (`mosquitto`)
- measurement database (`influxdb`)
- dashboard (`grafana`)
- plantor-backend, which glues all of these together

With additionally running [plantor](https://github.com/OliLay/plantor) on your microcontroller, 
you can view your plants health on a Grafana dashboard.

## Usage
Clone this repo and then execute  
 
`sudo sh setup.sh`  

`docker-compose up`

This will start all services, per default on the following ports:
- mosquitto `1883`
- influxdb `8086`
- grafana `3000`  

To open up your plant dashboard, simply type `localhost:3000` in your browser and
login with the default user/password `root/root`.