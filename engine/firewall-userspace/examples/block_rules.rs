//use futures::stream::StreamExt;
use redbpf::load::{Loader};
use std::env;
use std::process;
//use std::sync::{Arc, Mutex};
use tokio;
use tokio::runtime::Runtime;
use tokio::signal;
//use tokio::time::delay_for;
use redbpf::xdp;
use redbpf::{HashMap as BPFHashMap};



fn main() {
    if unsafe { libc::getuid() } != 0 { // Test if you are root
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }
    let args: Vec<String> = env::args().collect();

    let port :u16= args[1].parse::<u16>().unwrap();
    println!("Port {} is blocked",args[1]);


    let _ = Runtime::new().unwrap().block_on(async {

        
        let mut module = Loader::load(probe_code()).expect("error loading BPF program");
            for uprobe in module.xdps_mut() {
                uprobe.attach_xdp("wlp2s0", xdp::Flags::default()).unwrap();
            }
                
            //HashMap the

            let ports = BPFHashMap::<u16, u16>::new(module.map("port_map").unwrap()).unwrap();
            let key:u16=1;
            ports.set(key,port);

        signal::ctrl_c().await
    });

}

fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        env!("OUT_DIR"),
        "/target/bpf/programs/block_rules/block_rules.elf"
    ))
}