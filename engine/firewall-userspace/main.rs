#![no_std]
#![no_main]

use redbpf_probes::xdp::prelude::*;
use redbpf_probes::xdp::{PerfMap, XdpAction, XdpContext, MapData};
use redbpf_probes::kprobe::prelude::*;

program!(0xFFFFFFFE, "GPL");

#[kprobe("vfs_read")]
fn vfs_read_enter(_regs: Registers) {
    
}