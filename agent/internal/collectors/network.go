package collectors

import (
	"agent/internal/models"
	"fmt"
	"github.com/shirou/gopsutil/net"
)

// NetworkCollector собирает информацию о TCP и UDP портах
type NetworkCollector struct{}

func NewNetworkCollector() *NetworkCollector {
	return &NetworkCollector{}
}

// func (c *NetworkCollector) Collect(metrics *models.AgentMetrics) error {
// 	connections, err := net.Connections("all")
// 	if err != nil {
// 		return err
// 	}

// 	// Используем map для уникальности портов
// 	portMap := make(map[string]models.PortInfo)

// 	for _, conn := range connections {
// 		// Нас интересуют только прослушиваемые порты
// 		if conn.Status == "LISTEN" {
// 			port := conn.Laddr.Port
// 			protocol := "TCP"
// 			if conn.Type == "udp" {
// 				protocol = "UDP"
// 			}

// 			key := fmt.Sprintf("%s-%d", protocol, port)
// 			if _, exists := portMap[key]; !exists {
// 				portMap[key] = models.PortInfo{
// 					Port:     port,
// 					Protocol: protocol,
// 					State:    conn.Status,
// 				}
// 			}
// 		}
// 	}

// 	// Преобразуем map в slice
// 	var ports []models.PortInfo
// 	for _, port := range portMap {
// 		ports = append(ports, port)
// 	}

// 	metrics.Ports = ports
// 	return nil
// }

func (c *NetworkCollector) ChangeConfig(collType CollectorType, newconfig []string) {
	if collType == Network {
	}
}

func (c *NetworkCollector) Collect(metrics *models.AgentMetrics) error {
	// Собираем TCP соединения
	tcpConns, err := net.Connections("tcp")
	if err != nil {
		return err
	}

	// Собираем UDP соединения
	udpConns, err := net.Connections("udp")
	if err != nil {
		return err
	}

	connections := append(tcpConns, udpConns...)
	portMap := make(map[string]models.PortInfo)

	for _, conn := range connections {
		// Нас интересуют только активные соединения
		if conn.Status != "LISTEN" && conn.Status != "ESTABLISHED" {
			continue
		}

		port := conn.Laddr.Port
		protocol := "TCP"
		if conn.Type == 2 {
			protocol = "UDP"
		}

		key := fmt.Sprintf("%s-%d", protocol, port)
		if _, exists := portMap[key]; !exists {
			portMap[key] = models.PortInfo{
				Port:     uint16(port),
				Protocol: protocol,
				State:    conn.Status,
			}
		}
	}

	// Преобразуем map в slice
	var ports []models.PortInfo
	for _, port := range portMap {
		ports = append(ports, port)
	}

	metrics.Ports = ports
	return nil
}
