#![no_std]
#![no_main]

use redbpf_probes::tc::*;
use redbpf_probes::tc::prelude::*;
use core::mem::size_of;
use memoffset::offset_of;


program!(0xFFFFFFFE, "GPL");

#[tc_action]
pub fn block_port_80(skb: SkBuff) -> TcActionResult {

    let eth_proto: u16 = skb.load(offset_of!(ethhdr, h_proto))?;
    //Only look at IPv4 TCP packets
    if eth_proto as u32 != ETH_P_IP {
        return Ok(TcAction::Ok);
    }

    //get information
    let ip_start = size_of::<ethhdr>();
    let ip_proto: u8 = skb.load(ip_start + offset_of!(iphdr, protocol))?;
    let ip_len = ((skb.load::<u8>(ip_start)? & 0x0F) << 2) as usize;
    // Only look at TCP packets
    if ip_proto as u32 != IPPROTO_TCP {
        return Ok(TcAction::Ok);
    }

    let tcp_start = ip_start + ip_len;
    let dest_port: u16 = skb.load(tcp_start + offset_of!(tcphdr, dest))?;

    //if the port is not 80 pakcet pass
    if dest_port != 80 {
        return Ok(TcAction::Shot);
    }

    return Ok(TcAction::Shot);
}