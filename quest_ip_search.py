from scapy.all import ARP, Ether, srp




def scan(ip):
    arp_request = ARP(pdst=ip)
    ether = Ether(dst="ff:ff:ff:ff:ff:ff")
    packet = ether/arp_request

    result = srp(packet, timeout=1, verbose=0)[0]

    devices = {}
    for sent, received in result:
        devices[received.hwsrc] = received.psrc

    return devices




def search_oculus(devices:dict):
    prefixes = ["00:01:61", "80:f3:ef", "88:25:08", "94:f9:29", "b4:17:a8", "c0:dd:8a", "cc:a1:74",]
    oculuses = {}

    for key, value in devices.items():
        for prefix in prefixes:
            if key.startswith(prefix):
                oculuses[key] = value

    if not any(oculuses):
        oculuses = None

    return oculuses




if __name__ == "__main__":
    target_ip = "192.168.0.1/24"
    devices = scan(target_ip)
    print("===========================================")
    print(f"{target_ip} の範囲内で検出された機器\n")
    for k, v in devices.items():
        print(k, " | ", v)

    oculuses = search_oculus(devices)
    print("\n\n===========================================")
    print(f"{target_ip} の範囲内での QUEST っぽい奴\n")
    if oculuses is None:
        print("なし")
    else:
        for k, v in oculuses.items():
            print(k, " | ", v)

    input("\n\nEnterで終了")