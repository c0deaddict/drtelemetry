package telemetry

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
)

/**
 * https://docs.google.com/spreadsheets/d/1Yxw6iK_GtF4_63I4KzyF-M9t3f_XDMZrCgO7H-0lN0E/edit#gid=1084387652
 */
type TelemetryData struct {
	Time                   float32 `json:"time"`
	LapTime                float32 `json:"lapTime"`
	LapDistance            float32 `json:"lapDistance"`
	TotalDistance          float32 `json:"-"`
	X                      float32 `json:"-"` // World space position
	Y                      float32 `json:"-"` // World space position
	Z                      float32 `json:"-"` // World space position
	Speed                  float32 `json:"speed"`
	Xv                     float32 `json:"-"` // Velocity in world space
	Yv                     float32 `json:"-"` // Velocity in world space
	Zv                     float32 `json:"-"` // Velocity in world space
	Xr                     float32 `json:"-"` // rollx
	Yr                     float32 `json:"-"` // rolly
	Zr                     float32 `json:"-"` // rollz
	Xd                     float32 `json:"-"` // pitchx
	Yd                     float32 `json:"-"` // pitchy
	Zd                     float32 `json:"-"` // pitchz
	Susp_pos_bl            float32 `json:"suspensionPositionBL"`
	Susp_pos_br            float32 `json:"suspensionPositionBR"`
	Susp_pos_fl            float32 `json:"suspensionPositionFL"`
	Susp_pos_fr            float32 `json:"suspensionPositionFR"`
	Susp_vel_bl            float32 `json:"suspensionVelocityBL"`
	Susp_vel_br            float32 `json:"suspensionVelocityBR"`
	Susp_vel_fl            float32 `json:"suspensionVelocityFL"`
	Susp_vel_fr            float32 `json:"suspensionVelocityFR"`
	Wheel_speed_bl         float32 `json:"wheelSpeedBL"`
	Wheel_speed_br         float32 `json:"wheelSpeedBR"`
	Wheel_speed_fl         float32 `json:"wheelSpeedFL"`
	Wheel_speed_fr         float32 `json:"wheelSpeedFR"`
	Throttle               float32 `json:"throttlePosition"`
	Steer                  float32 `json:"steerPosition"`
	Brake                  float32 `json:"brakePosition"`
	Clutch                 float32 `json:"clutchPosition"`
	Gear                   float32 `json:"gear"`
	Gforce_lat             float32 `json:"gforceLat"`
	Gforce_lon             float32 `json:"gforceLon"`
	Lap                    float32 `json:"lap"`
	EngineRate             float32 `json:"rpm"`
	Sli_pro_native_support float32 `json:"-"`            // SLI Pro support
	Car_position           float32 `json:"-"`            // car race position
	Kers_level             float32 `json:"-"`            // kers energy left
	Kers_max_level         float32 `json:"-"`            // kers maximum energy
	Drs                    float32 `json:"-"`            // 0 = off, 1 = on
	Traction_control       float32 `json:"-"`            // 0 (off) - 2 (high)
	Anti_lock_brakes       float32 `json:"-"`            // 0 (off) - 1 (on)
	Fuel_in_tank           float32 `json:"-"`            // current fuel mass
	Fuel_capacity          float32 `json:"-"`            // fuel capacity
	In_pits                float32 `json:"-"`            // 0 = none, 1 = pitting, 2 = in pit area
	Sector                 float32 `json:"-"`            // 0 = sector1, 1 = sector2 float32 `json:"-"` 2 = sector3
	Sector1_time           float32 `json:"-"`            // time of sector1 (or 0)
	Sector2_time           float32 `json:"-"`            // time of sector2 (or 0)
	Brakes_temp_rl         float32 `json:"breaksTempRL"` // brakes temperature (centigrade)
	Brakes_temp_rr         float32 `json:"breaksTempRR"` // brakes temperature (centigrade)
	Brakes_temp_fl         float32 `json:"breaksTempFL"` // brakes temperature (centigrade)
	Brakes_temp_fr         float32 `json:"breaksTempFR"` // brakes temperature (centigrade)
	Tyre_pressure_rl       float32 `json:"-"`
	Tyre_pressure_rr       float32 `json:"-"`
	Tyre_pressure_fl       float32 `json:"-"`
	Tyre_pressure_fr       float32 `json:"-"`
	Laps_completed         float32 `json:"-"`
	Total_laps             float32 `json:"-"`
	Track_length           float32 `json:"trackLength"`
	Last_lap_time          float32 `json:"-"`        // last lap time
	Max_rpm                float32 `json:"maxRpm"`   // cars max RPM, at which point the rev limiter will kick in
	Idle_rpm               float32 `json:"idleRpm"`  // cars idle RPM
	Max_gears              float32 `json:"maxGears"` // maximum number of gears
}

var Addr = flag.String("udp", "localhost:20777", "udp service address")
var dataChannel = make(chan TelemetryData)

func ReadFromBytes(data []byte, telemetryData *TelemetryData) error {
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.LittleEndian, telemetryData); err != nil {
		return err
	}
	return nil
}

func RunServer() (chan TelemetryData, chan struct{}) {
	quit := make(chan struct{})

	go func() {
		protocol := "udp"

		udpAddr, err := net.ResolveUDPAddr(protocol, *Addr)
		if err != nil {
			fmt.Println("Wrong Address")
			return
		}

		udpConn, err := net.ListenUDP(protocol, udpAddr)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Listening for telemetry data on %s\n", udpAddr)
		defer udpConn.Close()

		for {
			select {
			default:
				HandleConnection(udpConn)
			case <-quit:
				fmt.Println("Stopping telemetry server")
				return
			}
		}
	}()

	return dataChannel, quit
}

func HandleConnection(conn *net.UDPConn) {
	buf := make([]byte, 264)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		fmt.Println("Error Reading")
		fmt.Println(err)
		return
	} else {
		var reader = bytes.NewBuffer(buf[:n])
		var foo TelemetryData
		err := binary.Read(reader, binary.LittleEndian, &foo)
		if err != nil {
			fmt.Println(err)
		}
		dataChannel <- foo
	}
}
