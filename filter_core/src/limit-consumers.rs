#![no_std]
#![no_main]
use core::mem::size_of;
use memoffset::offset_of;
use core::marker::PhantomData;
use redbpf_macros::map;
use redbpf_probes::tc::prelude::*;
use redbpf_probes::tc::{TcAction, TcActionResult};

program!(0xFFFFFFFE, "GPL");

const PIN_GLOBAL_NS: u32 = 2;
#[repr(C)]
struct bpf_elf_map {
    type_: u32,
    size_key: u32,
    size_value: u32,
    max_elem: u32,
    flags: u32,
    id: u32,
    pinning: u32,
}

pub struct TcHashMap<K, V> {
    def: bpf_elf_map,
    _k: PhantomData<K>,
    _v: PhantomData<V>,
}

impl<K, V> TcHashMap<K, V> {
    /// Creates a map with the specified maximum number of elements.
    pub const fn with_max_entries(max_entries: u32) -> Self {
        Self {
            def: bpf_elf_map {
                type_: 1, // BPF_MAP_TYPE_HASH
                size_key: mem::size_of::<K>() as u32,
                size_value: mem::size_of::<V>() as u32,
                max_elem: max_entries,
                flags: 0,
                id: 0,
                pinning: PIN_GLOBAL_NS,
            },
            _k: PhantomData,
            _v: PhantomData,
        }
    }
    /// Returns a reference to the value corresponding to the key.
    #[inline]
    pub fn get(&mut self, key: &K) -> Option<&V> {
        unsafe {
            let value = bpf_map_lookup_elem(
                &mut self.def as *mut _ as *mut _,
                key as *const _ as *const _,
            );
            if value.is_null() {
                None
            } else {
                Some(&*(value as *const V))
            }
        }
    }

    #[inline]
    pub fn get_mut(&mut self, key: &K) -> Option<&mut V> {
        unsafe {
            let value = bpf_map_lookup_elem(
                &mut self.def as *mut _ as *mut _,
                key as *const _ as *const _,
            );
            if value.is_null() {
                None
            } else {
                Some(&mut *(value as *mut V))
            }
        }
    }

    /// Set the `value` in the map for `key`
    #[inline]
    pub fn set(&mut self, key: &K, value: &V) {
        unsafe {
            bpf_map_update_elem(
                &mut self.def as *mut _ as *mut _,
                key as *const _ as *const _,
                value as *const _ as *const _,
                BPF_ANY.into(),
            );
        }
    }
   /// Delete the entry indexed by `key`
   #[inline]
   pub fn delete(&mut self, key: &K) {
       unsafe {
           bpf_map_delete_elem(
               &mut self.def as *mut _ as *mut _,
               key as *const _ as *const _,
           );
       }
   }
}



#[derive(Clone, Debug)]
#[repr(C)]
struct Source {
  addr: u32,
  port: u32, // should be u16, but need padding (?)
}
//static mut blocked_packets: HashMap<Source, u8> = HashMap::with_max_entries(10240);

#[map(link_section = "maps")]
static mut blocked_packets: TcHashMap<u64, u64> = TcHashMap::<u64, u64>::with_max_entries(1024);


fn is_port_in_map(src_port: u16,src_addr: u32 ) -> bool{
  let src = Source {
    addr: src_addr,
    port: src_port as u32,
  };
  let blockedPaq: Option<&mut u8>;
  unsafe {
  //blockedPaq= blocked_packets.get_mut(&src);
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

  /** 
  //Lookup if we have a filter rule and delete it
  if is_port_in_map(src_port, src_addr){
    unsafe {blocked_packets.delete(&src);}
    return Ok(TcAction::Shot);
  }**/

  unsafe {
    let key = 0;
    if let Some(mut cnt) = blocked_packets.get_mut(&key) {
        *cnt += 1;
    } else {
        let val = 0;
        blocked_packets.set(&key, &val);
    }
}
  

  Ok(TcAction::Ok)
}
