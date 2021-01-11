use redbpf::load::{ Loader};
use std::env;
use std::process;
use std::time::Duration;
use tokio;
use tokio::runtime::Runtime;
use tokio::signal;
use tokio::time::delay_for;
use redbpf::xdp;
use redbpf::HashMap as BPFHashMap;
use simple_logger::SimpleLogger;
use userspace_map::IPAggs;
use userspace_map::PortAggs;

pub mod userspace_map;
pub mod network_tool;
pub mod elastic_mapping;

fn main() {
    if unsafe { libc::getuid() } != 0 {
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }

    SimpleLogger::new()
    .with_level(log::LevelFilter::Debug)
    .init()
    .unwrap();
   
    let _ = Runtime::new().unwrap().block_on(async {

        
        let mut module = Loader::load(probe_code()).expect("error loading BPF program");


            for uprobe in module.xdps_mut() {
                uprobe.attach_xdp("wlp2s0", xdp::Flags::default()).unwrap();
            }
           log::debug!("ave");


        tokio::spawn(async move {
            let _ips = BPFHashMap::<u32, IPAggs>::new(module.map("ip_map").unwrap()).unwrap();
            let _ports = BPFHashMap::<u16, PortAggs>::new(module.map("port_map").unwrap()).unwrap();
            loop {
                delay_for(Duration::from_millis(60000)).await;
                //format ips Hashmap into vec
                let ip_vec: Vec<(u32, IPAggs)> = _ips.iter().collect();
                let mut parsed_ips: Vec<elastic_mapping::EsReadyIpAggs> = Vec::new();

                log::debug!("========ip addresses=======");

                for (k, v) in ip_vec.iter().rev() {
                    log::debug!(
                        "{:?} - > count:{:?}",
                        network_tool::u32_to_ipv4(*k),
                        v.count
                    );

                    let current_ip_agg = elastic_mapping::EsReadyIpAggs {
                        ip: network_tool::u32_to_ipv4(*k).to_string(),
                        count: v.count,
                        usage: v.usage,
                    };

                    parsed_ips.push(current_ip_agg);
                    _ips.delete(*k);
                }

            }
        });

        signal::ctrl_c().await
    });
}

fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        env!("OUT_DIR"),
        "/target/bpf/programs/block_rules/block_rules.elf"
    ))
}
