package app

import (
	"fmt"
	"log/slog"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/songgao/water"
)

type TUN struct {
	logger *slog.Logger
}

func NewTUN(logger *slog.Logger) *TUN {
	return &TUN{logger: logger}
}

func (t *TUN) MustRun(showByteLimit int) {
	if err := t.run(showByteLimit); err != nil{
		panic("TUN interface creation failed")
	}
}

func (t *TUN) run(showByteLimit int) error {

	// Interface creation
	config := water.Config{
        DeviceType: water.TUN,
    }
    iface, err := water.New(config)
    if err != nil {
        return err
    }

    defer iface.Close()

    // Package buffer (MTU 1500)
    packet := make([]byte, 1500)

    for {
        // Package reading from TUN
        n, err := iface.Read(packet)
        if err != nil {
			t.logger.Debug("Package reading error", "error", err.Error())
			break
        }

        // gopacket lib -- Parsing package
        parsedPacket := gopacket.NewPacket(
            packet[:n], 
            layers.LayerTypeIPv4, 
            gopacket.DecodeOptions{NoCopy: true, Lazy: true},
        )

        // Getting ipv4 header
        ipLayer := parsedPacket.Layer(layers.LayerTypeIPv4)
        if ipLayer == nil {
            continue // if there is no ipv4
        }
        ip, _ := ipLayer.(*layers.IPv4)

        // Getting transport protocol + from/to ports

        var transportProtocol string
        var tcp *layers.TCP
        var udp *layers.UDP
        var senderPort string
        var receiverPort string
        switch ip.Protocol {
        case layers.IPProtocolTCP:
            tcpLayer := parsedPacket.Layer(layers.LayerTypeTCP)
            if tcpLayer != nil {
                tcp, _ = tcpLayer.(*layers.TCP)
                transportProtocol = "TCP"
                senderPort = fmt.Sprintf("%d", tcp.SrcPort)
                receiverPort = fmt.Sprintf("%d", tcp.DstPort)
            }
        case layers.IPProtocolUDP:
            udpLayer := parsedPacket.Layer(layers.LayerTypeUDP)
            if udpLayer != nil {
                udp, _ = udpLayer.(*layers.UDP)
                transportProtocol = "UDP"
                senderPort = fmt.Sprintf("%d", udp.SrcPort)
                receiverPort = fmt.Sprintf("%d", udp.DstPort)
            }
        default:
            transportProtocol = fmt.Sprintf("%s, (code: %d)", ip.Protocol, ip.Protocol)
        }

        var packetData []byte
		if n <= showByteLimit {
            packetData = packet[:n]
		} else {
            packetData = packet[:showByteLimit]
		}

        // Content without headers
        var payloadContent string
		payload := parsedPacket.ApplicationLayer()
        if payload != nil {
            payloadContent = fmt.Sprintf("%s\n", payload.Payload())
        }

        // logging
        t.logger.Debug(fmt.Sprintf("Listen packages on TUN interface: %s", iface.Name()),
            "IP Sender", fmt.Sprintf("%s",ip.SrcIP),
            "IP Receiver", fmt.Sprintf("%s",ip.DstIP),
            "Transport protoocol", transportProtocol,
            "Sender port", senderPort,
            "Receiver port", receiverPort,
            "Package content", fmt.Sprintf("%b", packetData),
            "Payload", payloadContent,
        )

    }
	
	return nil
}