use redbpf::load::{ Loader};
use std::env;
use std::process;
use std::time::Duration;
use std::ptr;
use std::sync::{Arc, Mutex};
use tokio;
use tokio::runtime::Runtime;
use tokio::signal;
use tokio::time::delay_for;
use redbpf::Module;
use redbpf::xdp;
use redbpf::HashMap as BPFHashMap;
use probes::block_rules::PortAggs;
use probes::block_rules::IPAggs;



fn main() {
    if unsafe { libc::getuid() } != 0 {
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }
   
    let _ = Runtime::new().unwrap().block_on(async {

        
        let mut module = Loader::load(probe_code()).expect("error loading BPF program");
        

        let ips = BPFHashMap::<u32, IPAggs>::new(module.map("ip_map").unwrap()).unwrap();
        let ports = BPFHashMap::<u16, PortAggs>::new(module.map("port_map").unwrap()).unwrap();

            for uprobe in module.xdps_mut() {
                uprobe.attach_xdp("wlp2s0", xdp::Flags::default()).unwrap();
            }

           

            

        signal::ctrl_c().await
    });
}

fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        env!("OUT_DIR"),
        "/target/bpf/programs/block_rules/block_rules.elf"
    ))
}
