import subprocess, re
import ping3




print("\n========== QUEST IP SEARCHER ===========\n")




def scan_devices(start:int, end:int) -> list:
    print(f"範囲 : 192.168.0.{start} ~ 192.168.0.{end-1}")

    devices = []
    for _ in range(start, end):
        scan_ip = f"192.168.0.{_}"
        progress = round(_ * 0.4)
        if progress >= 100:
            progress = 100
        progress_bar = ("#" * round(progress * 0.4)).ljust(40, "-")
        print("\r" + "スキャン中 : ", scan_ip , f"\n{progress_bar} {str(progress)}%" + "\033[1A", end="")

        result = ping3.ping(scan_ip, timeout=0.2)

        if not result in [False, None]:
            devices.append(scan_ip)

    return devices




def get_macaddr(ips:list):
    mac_ip = {}

    subprocess.run(["chcp", "65001"], shell=True, text=True, stdout=subprocess.PIPE)

    for ip in ips:
        result = subprocess.check_output(["arp", "-a", ip], text=True)

        mac_address_match = re.compile(r"([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})").search(result)
        if mac_address_match:
            mac_ip[mac_address_match.group()] = ip
        else:
            mac_ip[f"{ip}_not_found_mac"] = ip

    return mac_ip




def search_oculus(devices:dict):
    prefixes = ["00-01-61", "80-f3-ef", "88-25-08", "94-f9-29", "b4-17-a8", "c0-dd-8a", "cc-a1-74"]
    oculuses = {}

    for key, value in devices.items():
        for prefix in prefixes:
            if key.startswith(prefix):
                oculuses[key] = value

    if not any(oculuses):
        oculuses = None

    return oculuses




if __name__ == '__main__':
    existing_ip = scan_devices(1, 255)
    devices = get_macaddr(existing_ip)

    print(f"\n\n\n\n========== 検出された機器 ==============\n")
    for k, v in devices.items():
        print(k, " | ", v)

    print(f"\n\n\n========== QUESTっぽい機器 =============\n")
    quests = search_oculus(devices)
    if not quests:
        print("なし")
    else:
        for k, v in quests.items():
            print(k, " | ", v)
    print("\n\n\n")
    input("Enterで終了")