#!/usr/bin/python3

import sys
import platform
import subprocess
from multiprocessing import Process
from time import sleep


def ping(ipaddr):
    param = 'n' if platform.system().lower() == 'windows' else '-c'
    command = ['ping', param, '1', ipaddr]
    ping_output = subprocess.call(command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    if ping_output == 0:
        print(ipaddr)


def get_ip_range(octets):
    if octets[0] == 192 and octets[1] == 168:
        addr = {"oct1": range(192, 193), "oct2": range(168, 169),
                "oct3": range(0, 256), "oct4": range(0, 256)}
    elif octets[0] == 172 and 16 <= octets[1] <= 31:
        addr = {"oct1": range(172, 173), "oct2": range(16, 32),
                "oct3": range(0, 256), "oct4": range(0, 256)}
    elif octets[0] == 10:
        addr = {"oct1": range(10, 11), "oct2": range(0, 256),
                "oct3": range(0, 256), "oct4": range(0, 256)}
    else:
        print("IP Address is not in a private range.")
        sys.exit(1)

    threads = 0
    max_threads = 1000
    for oct1 in addr['oct1']:
        for oct2 in addr['oct2']:
            for oct3 in addr['oct3']:
                for oct4 in addr['oct4']:
                    ipaddr = str(oct1) + "." + str(oct2) + "." + str(oct3) + "." + str(oct4)
                    if threads > max_threads:
                        sleep(5)
                        threads = 0
                    else:
                        threads += 1
                    Process(target=ping,
                            args=(ipaddr,)).start()


def main():
    if len(sys.argv) != 2:
        print("Syntax: ./ipsweep.py <ip address>")
        sys.exit(1)
    else:
        octets = list(map(int, sys.argv[1].split(".")))

    get_ip_range(octets)


main()
