use redbpf::{load::Loaded, xdp::{self}};
use redbpf::{load::Loader, HashMap as BPFHashMap};
use std::process;
use std::time::Duration;
use tokio::{runtime::Runtime, signal};
use tokio::time::delay_for;
use crate::aggs::RequestInfo;

pub mod aggs;
pub mod elastic_mapping;
pub mod network_utils;


fn start_perf_event_handler(loader: Loaded) {
    
    // Listen to incoming map's data
    tokio::spawn(async move {
        // Load the Hashmap into variables
        let ips = BPFHashMap::<u32, aggs::IPAggs>::new(loader.map("ip_map").unwrap()).unwrap();
        let ports = BPFHashMap::<u16, aggs::PortAggs>::new(loader.map("port_map").unwrap()).unwrap();
        //let datatab = BPFHashMap::<u32, RequestInfo>::new(loader.map("requests").unwrap()).unwrap();

        loop {
            // Wait for ?s
            delay_for(Duration::from_millis(10000)).await;
            
            // Format data Hashmap into vec
            /*
            let data_vec: Vec<(u32, RequestInfo)> = datatab.iter().collect();

            println!("========packet data=======");
            
            for (k, v) in data_vec.iter().rev() {
                
                println!(
                    "{:?} - > data:{:?}",
                    k,
                    v,
                );
                
            }*/

            // Format ips Hashmap into vec
            let ip_vec: Vec<(u32, aggs::IPAggs)> = ips.iter().collect();
            let mut parsed_ips: Vec<elastic_mapping::EsReadyIpAggs> = Vec::new();

            println!("========ip addresses=======");
            
            for (k, v) in ip_vec.iter().rev() {
                println!(
                    "{:?} - > count:{:?}",
                    network_utils::u32_to_ipv4(*k),
                    v.count,
                );

                let current_ip_agg = elastic_mapping::EsReadyIpAggs {
                    ip: network_utils::u32_to_ipv4(*k).to_string(),
                    count: v.count,
                    usage: v.usage,
                };

                parsed_ips.push(current_ip_agg);
                ips.delete(*k);
            }

            // Format port Hashmap into vec
            let port_vec: Vec<(u16, aggs::PortAggs)> = ports.iter().collect();
            let mut parsed_ports: Vec<elastic_mapping::EsReadyPortAggs> = Vec::new();

            println!("========ports=======");

            for (k, v) in port_vec.iter().rev() {
                println!("{:?} - > count:{:?}", k, v.count);

                let current_port_agg = elastic_mapping::EsReadyPortAggs {
                    port: *k,
                    count: v.count,
                };

                parsed_ports.push(current_port_agg);
                ports.delete(*k);
            }

            
            let data_iteration = elastic_mapping::BPFDataIteration {
                ips: parsed_ips,
                ports: parsed_ports,
            };

            println!("========data========");

            println!("IPS:{:?}", data_iteration.ips);
            println!("PORTS:{:?}", data_iteration.ports);
            
        }
    });
}


fn main() {
    if unsafe { libc::getuid() } != 0 {
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }

    let _ = Runtime::new().unwrap().block_on(async {

        let mut loader = Loader::load(probe_code()).expect("error loading probe");
        
        /*
        for uprobe in loader.xdps_mut() {
            uprobe.attach_xdp("wlp2s0", xdp::Flags::default()).unwrap();
        }
        */

        for kp in loader.kprobes_mut() {
            kp.attach_kprobe(&kp.name(), 0)
                .expect(&format!("error attaching kprobe program {}", kp.name()));
        }

        println!("========Success to attach ========");

        //start_perf_event_handler(loader);

        signal::ctrl_c().await
    });
}


fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        "/home/matthieu/hello-bpf/target/bpf/programs/block_http/block_http.elf"
    ))
}
