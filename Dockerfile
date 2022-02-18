FROM eclipse-mosquitto:2.0.14
RUN mkdir -p /mosquitto/data/
COPY mosquitto.conf /mosquitto/config/mosquitto.conf