use std::fs::File;
use std::collections::HashMap;
use std::io::{self, BufRead};
use std::path::Path;
use redbpf::load::{Loader};
use std::env;
use std::process;
use tokio;
use tokio::signal;
use redbpf::xdp;
use redbpf::{HashMap as BPFHashMap};
use tokio::time::delay_for;
use std::time::{Duration};

fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}

fn read_port(port:String) -> Result<u16, &'static str>{
    match port.parse::<u16>(){
        Ok(port) => return Ok(port),
        Err(_err)=> return Err(""),
    }
}



fn read_blocked_ports() -> HashMap<u16, u16>{
    let mut ports= HashMap::new();
    println!("Read the file : ");
    if let Ok(lines) = read_lines("./blockedPort") { // Read the file blocked port
        // Consumes the iterator, returns an (Optional) String
        for line in lines {
            if let Ok(port_read) = line {  // Read the port line by line
                match read_port(port_read.clone()) {  // Test if the port on the ligne is a u16
                    Ok(port) =>{
                        println!("      Port {} blocked", port);
                        ports.insert(port,1);
                    }
                    Err(_err)=>eprintln!("      {} is not a port ",port_read),
                }

                 
                
            }
        }
    }
    return ports;
}

#[tokio::main]
async fn main() -> Result<(), io::Error> {

    
    
    if unsafe { libc::getuid() } != 0 { // Test if you are root
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }

   
        
        let mut module = Loader::load(probe_code()).expect("error loading BPF program");
            for uprobe in module.xdps_mut() {
                uprobe.attach_xdp("wlp2s0", xdp::Flags::default()).unwrap(); // Attach the ELF bytecode to the wifi XDP interface
            }
                
            
            
            tokio::spawn(async move { // Open a thread to read the file while the progamm still running
                let mut map_file_reader :HashMap<u16, u16>;
                let ports = BPFHashMap::<u16, u16>::new(module.map("port_map").unwrap()).unwrap(); // Create the HashMap to send data to the user space
                loop{
                    
                    map_file_reader=read_blocked_ports();
                    //UPDATE the BPF map
                    //Delete ports you are not in the file 
                    for (k,_v) in ports.iter(){
                        match map_file_reader.get(&k){
                            Some(_)=>(), 
                            None => ports.delete(k),
                        }
                    }
                    //Add the new ports
                    for (k,_v) in map_file_reader.iter(){
                        match ports.get(*k){
                            Some(_)=>(), 
                            None => ports.set(*k,1),
                        }
                    }
                    delay_for(Duration::from_millis(3000)).await;
                }
            });   

            

        signal::ctrl_c().await
    

}

fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        env!("OUT_DIR"),
        "/target/bpf/programs/block_rules/block_rules.elf"
    ))
}