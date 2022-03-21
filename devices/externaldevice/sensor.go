package ExternalDevice

import (
	"encoding/json"
	common "github.com/W-Floyd/ha-mqtt-iot/common"
	"log"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// Do not modify this file, it is automatically generated
////////////////////////////////////////////////////////////////////////////////
func (d Sensor) GetRawId() string {
	return "sensor"
}
func (d Sensor) AddMessageHandler() {
	d.MQTT.MessageHandler = MakeMessageHandler(d)
}
func (d Sensor) GetUniqueId() string {
	return d.UniqueId
}
func (d Sensor) PopulateDevice() {
	d.Device.Manufacturer = Manufacturer
	d.Device.Model = SoftwareName
	d.Device.Name = InstanceName
	d.Device.SwVersion = SWVersion
}

type Sensor struct {
	AvailabilityMode     string `json:"availability_mode"`
	AvailabilityTemplate string `json:"availability_template"`
	AvailabilityTopic    string `json:"availability_topic"`
	Device               struct {
		ConfigurationUrl string   `json:"configuration_url"`
		Connections      []string `json:"connections"`
		Identifiers      []string `json:"identifiers"`
		Manufacturer     string   `json:"manufacturer"`
		Model            string   `json:"model"`
		Name             string   `json:"name"`
		SuggestedArea    string   `json:"suggested_area"`
		SwVersion        string   `json:"sw_version"`
		Viadevice        string   `json:"viadevice"`
	} `json:"device"`
	DeviceClass            string        `json:"device_class"`
	EnabledByDefault       bool          `json:"enabled_by_default"`
	Encoding               string        `json:"encoding"`
	EntityCategory         string        `json:"entity_category"`
	ExpireAfter            int           `json:"expire_after"`
	ForceUpdate            bool          `json:"force_update"`
	Icon                   string        `json:"icon"`
	LastResetValueTemplate string        `json:"last_reset_value_template"`
	Name                   string        `json:"name"`
	ObjectId               string        `json:"object_id"`
	PayloadAvailable       string        `json:"payload_available"`
	PayloadNotAvailable    string        `json:"payload_not_available"`
	Qos                    int           `json:"qos"`
	StateClass             string        `json:"state_class"`
	StateTopic             string        `json:"state_topic"`
	StateFunc              func() string `json:"-"`
	UniqueId               string        `json:"unique_id"`
	UnitOfMeasurement      string        `json:"unit_of_measurement"`
	ValueTemplate          string        `json:"value_template"`
	MQTT                   MQTTFields    `json:"-"`
}

func (d Sensor) UpdateState() {
	if d.StateTopic != "" {
		state := d.StateFunc()
		if state != stateStore.Sensor.State[d.UniqueId] || d.MQTT.ForceUpdate {
			c := *d.MQTT.Client
			token := c.Publish(d.StateTopic, common.QoS, common.Retain, state)
			stateStore.Sensor.State[d.UniqueId] = state
			token.Wait()
		}
	}
}
func (d Sensor) Subscribe() {
	c := *d.MQTT.Client
	message, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	token := c.Publish(GetDiscoveryTopic(d), 0, true, message)
	token.Wait()
	time.Sleep(common.HADiscoveryDelay)
	d.AnnounceAvailable()
	d.UpdateState()
}
func (d Sensor) UnSubscribe() {
	c := *d.MQTT.Client
	token := c.Publish(d.AvailabilityTopic, common.QoS, common.Retain, "offline")
	token.Wait()
}
func (d Sensor) AnnounceAvailable() {
	c := *d.MQTT.Client
	token := c.Publish(d.AvailabilityTopic, common.QoS, common.Retain, "online")
	token.Wait()
}
func (d Sensor) Initialize() {
	d.PopulateDevice()
	d.PopulateTopics()
	d.AddMessageHandler()
}
func (d Sensor) PopulateTopics() {
	if d.StateFunc != nil {
		d.StateTopic = GetTopic(d, "state_topic")
	}
}
