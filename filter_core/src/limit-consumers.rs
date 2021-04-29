#![no_std]
#![no_main]
use core::mem::size_of;
use memoffset::offset_of;
use redbpf_macros::map;
use redbpf_probes::tc::prelude::*;

program!(0xFFFFFFFE, "GPL");

#[derive(Clone, Debug)]
#[repr(C)]
struct Source {
  addr: u32,
  port: u32, // should be u16, but need padding (?)
}

#[map(link_section = "maps")]
static mut blocked_packets: HashMap<Source, u8> = HashMap::with_max_entries(10240);



fn is_port_in_map(src_port: u16,src_addr: u32 ) -> bool{
  let src = Source {
    addr: src_addr,
    port: src_port as u32,
  };
  let blockedPaq: Option<&mut u8>;
  unsafe {
  blockedPaq= blocked_packets.get_mut(&src);
 }
      match blockedPaq{
          // The port is in the map block it 
          Some(_x)    => return true,
          // The port is not in the map
          None    => return false,
      }
  
}

#[tc_action]
fn limit(skb: SkBuff) -> TcActionResult {
  let eth_proto: u16 = skb.load(offset_of!(ethhdr, h_proto))?;
  //Only look at IPv4 TCP packets
  if eth_proto as u32 != ETH_P_IP {
    return Ok(TcAction::Ok);
  }

  let ip_start = size_of::<ethhdr>();
  let ip_proto: u8 = skb.load(ip_start + offset_of!(iphdr, protocol))?;
  let ip_len = ((skb.load::<u8>(ip_start)? & 0x0F) << 2) as usize;
  // Only look at TCP packets
  if ip_proto as u32 != IPPROTO_TCP {
    return Ok(TcAction::Ok);
  }

  let tcp_start = ip_start + ip_len;
  //let dest_port: u16 = skb.load(tcp_start + offset_of!(tcphdr, dest))?;
  // Only look at The port we want
  /**if dest_port != 443 {
    return Ok(TcAction::Ok);
  }**/

  let data_offset = (skb.load::<u8>(tcp_start + 12)? >> 4) << 2;
  //let data_start = tcp_start + data_offset as usize;

  //assignation to a new struct 
  let src_addr: u32 = skb.load(ip_start + offset_of!(iphdr, saddr))?;
  let src_port: u16 = skb.load(tcp_start + offset_of!(tcphdr, source))?;
  let src = Source {
    addr: src_addr,
    port: src_port as u32,
  };

  //Lookup if we have a filter rule and delete it
  if is_port_in_map(src_port, src_addr){
    unsafe {blocked_packets.delete(&src);}
    return Ok(TcAction::Shot);
  }
  

  Ok(TcAction::Ok)
}
