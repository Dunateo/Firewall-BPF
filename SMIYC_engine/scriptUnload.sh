#!/bin/bash
Device="wlp2s0"

sudo tc filter delete dev $Device egress pref 1
sudo tc qdisc del dev $Device clsact
