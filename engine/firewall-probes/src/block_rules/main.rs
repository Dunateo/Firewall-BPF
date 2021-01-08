#![no_std]
#![no_main]

use redbpf_probes::xdp::prelude::*;

program!(0xFFFFFFFE, "GPL");

#[map("myMap")]
static mut myMap : HashMap<u64,u64> = HashMap::with_max_entries(10240);


// A macro which defines an insert behavior
// If there isn't a key, create an entry
// if there is a key, return the entry

#[map("port_map")]
static mut port_map: HashMap<u16, u16> = HashMap::with_max_entries(10240);

#[xdp]
pub fn block_port_80(ctx: XdpContext) -> XdpResult {
    let key:u16=1;
    let port :u16;
    unsafe{
        port =*port_map.get(&key).unwrap();
    }
    
    if let Ok(transport) = ctx.transport() {
        
        if transport.dest() == port {
            return Ok(XdpAction::Drop);
            
        }else if transport.source() == port  {
            return Ok(XdpAction::Drop);
        }
    }

    Ok(XdpAction::Pass)
}