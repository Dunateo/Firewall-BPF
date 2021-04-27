Device="wlp2s0"

sudo tc qdisc add dev $Device clsact
sudo tc filter add dev $Device egress pref 1 handle 1 bpf da obj target/bpf/programs/limit/limit.elf sec tc_action/limit
