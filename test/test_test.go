package test

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/lbsystem/protocol"
)

// prepareUDPEthernetFrame 构建一个包含UDP数据的以太网帧
func prepareUDPEthernetFrame() []byte {
	// 以太网帧头部：6字节目的MAC + 6字节源MAC + 2字节类型
	dstMAC := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} // 广播地址
	srcMAC := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06} // 示例地址
	etherType := []byte{0x08, 0x00}                      // 表示IPv4

	// IPv4头部：20字节
	ipHeader := make([]byte, 1400)
	ipHeader[0] = 0x45                                    // IPv4, 头部长度20字节
	ipHeader[1] = 0x00                                    // 服务类型
	binary.BigEndian.PutUint16(ipHeader[2:4], 40)         // 总长度：IP头 + UDP头 + 数据 = 20 + 8 + 12
	binary.BigEndian.PutUint16(ipHeader[4:6], 0)          // 标识
	binary.BigEndian.PutUint16(ipHeader[6:8], 0)          // 标志和片偏移
	ipHeader[8] = 64                                      // TTL
	ipHeader[9] = 17                                      // 协议：UDP
	binary.BigEndian.PutUint16(ipHeader[10:12], 0)        // 校验和，先设为0，可以实际计算
	copy(ipHeader[12:16], []byte{0x0A, 0x00, 0x00, 0x01}) // 源IP地址
	copy(ipHeader[16:20], []byte{0x0A, 0x00, 0x00, 0x02}) // 目标IP地址

	// UDP头部：8字节
	udpHeader := make([]byte, 8)
	binary.BigEndian.PutUint16(udpHeader[0:2], 12345) // 源端口
	binary.BigEndian.PutUint16(udpHeader[2:4], 12346) // 目标端口
	binary.BigEndian.PutUint16(udpHeader[4:6], 20)    // UDP长度：8头部 + 12数据
	binary.BigEndian.PutUint16(udpHeader[6:8], 0)     // 校验和，可以实际计算

	// UDP数据
	udpData := []byte("Hello, World") // 12字节数据

	// 组装完整的以太网帧
	frame := append(dstMAC, srcMAC...)
	frame = append(frame, etherType...)
	frame = append(frame, ipHeader...)
	frame = append(frame, udpHeader...)
	frame = append(frame, udpData...)

	return frame
}

// BenchmarkUnmarshalBinary 测试 Ethernet.UnmarshalBinary 方法的性能
func BenchmarkUnmarshalBinary(b *testing.B) {
	data := prepareUDPEthernetFrame()
	eth := protocol.NewEthernet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := eth.UnmarshalBinary(data); err != nil {
			b.Fatal(err)
		}
	}
}
func prepareIPv6UDPEthernetFrame() []byte {
	dstMAC := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	srcMAC := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	etherType := []byte{0x86, 0xDD} // IPv6 的以太网类型

	// IPv6头部
	ipv6Header := make([]byte, 40)
	ipv6Header[0] = 0x60                                 // 版本号 (4 bits) 和流量类别的前半部分 (4 bits)
	ipv6Header[1] = 0x00                                 // 流量类别的后半部分 (4 bits) 和流标签的前半部分 (4 bits)
	binary.BigEndian.PutUint32(ipv6Header[2:6], 0)       // 流标签的后半部分 (16 bits) 和载荷长度 (16 bits)
	ipv6Header[6] = 17                                   // 下一个头部：UDP
	ipv6Header[7] = 64                                   // 跳限制
	copy(ipv6Header[8:24], net.ParseIP("2001:0db8::1"))  // 源IPv6地址
	copy(ipv6Header[24:40], net.ParseIP("2001:0db8::2")) // 目标IPv6地址

	// UDP头部
	udpHeader := make([]byte, 8)
	binary.BigEndian.PutUint16(udpHeader[0:2], 12345) // 源端口
	binary.BigEndian.PutUint16(udpHeader[2:4], 12346) // 目标端口
	binary.BigEndian.PutUint16(udpHeader[4:6], 28)    // UDP长度：8头部 + 20数据
	binary.BigEndian.PutUint16(udpHeader[6:8], 0)     // 校验和，可以实际计算

	// UDP数据
	udpData := []byte("Hello, World from IPv6 UDP")

	// 组装完整的以太网帧
	frame := append(dstMAC, srcMAC...)
	frame = append(frame, etherType...)
	frame = append(frame, ipv6Header...)
	frame = append(frame, udpHeader...)
	frame = append(frame, udpData...)

	return frame
}

// BenchmarkUnmarshalIPv6UDP 测试解析IPv6 UDP Ethernet帧的性能
func BenchmarkUnmarshalIPv6UDP(b *testing.B) {
	data := prepareIPv6UDPEthernetFrame()
	eth := protocol.NewEthernet()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := eth.UnmarshalBinary(data); err != nil {
			b.Fatal(err)
		}
	}
}
