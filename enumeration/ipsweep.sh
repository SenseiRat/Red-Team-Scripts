#!/bin/bash

if [[ -z $1 ]]; then
    echo "Syntax: ./ipsweep.sh <ip address>"
    exit 1
else
    addr=( $(echo $1 | sed 's/\./ /g') )
    if [[ "${addr[0]}" -eq 192 ]] && [[ "${addr[1]}" -eq 168 ]]; then
        oct1=192
    elif [[ "${addr[0]}" -eq 172 ]]; then
        if [[ "${addr[1]}" -ge 16 ]] && [[ "${addr[1]}" -le 31 ]]; then
            oct1=172
        fi
    elif [[ "${addr[0]}" -eq 10 ]]; then
        oct1=10
    fi
fi

if [[ -z ${oct1} ]]; then
    echo "IP Address is not in a private range."
    exit 1
elif [[ ${oct1} -eq 172 ]]; then
    oct2_min=16
    oct2_max=31
elif [[ ${oct1} -eq 192 ]]; then
    oct2_min=168
    oct2_max=168
elif [[ ${oct1} -eq 10 ]]; then
    oct2_min=0
    oct2_max=255
fi

threads=1
for oct2 in $(seq ${oct2_min} ${oct2_max}); do
    for oct3 in $(seq 0 255); do
        for oct4 in $(seq 0 255); do
            if [[ ${threads} -eq 1000 ]]; then
                sleep 1
                threads=1
            fi
            ping -c 1 ${oct1}.${oct2}.${oct3}.${oct4} 2>/dev/null | grep "64 bytes" | sed 's/^.*from \(.*\):.*$/\1/' &
            threads=$(expr ${threads} + 1)
        done
    done
done

