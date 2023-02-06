package hass

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MessageHandler func(mqtt.Message)
