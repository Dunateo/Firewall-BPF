#![no_std]
#![no_main]

use redbpf_probes::xdp::prelude::*;

program!(0xFFFFFFFE, "GPL");


#[map("port_map")]
static mut port_map: HashMap<u16, u16> = HashMap::with_max_entries(10240); // Map used to get data from user space


fn is_port_in_map(port:u16) -> bool{
    let blockedPorts: Option<&mut u16>;
    unsafe {
    blockedPorts= port_map.get_mut(&port);
   }
        match blockedPorts{
            // The port is in the map block it 
            Some(_x)    => return true,
            // The port is not in the map
            None    => return false,
        }
    
}


#[xdp] // Attach the function to the XDP
pub fn block_port_80(ctx: XdpContext) -> XdpResult {
    
    
    if let Ok(transport) = ctx.transport() {
        
        if is_port_in_map(transport.dest() )  {
            return Ok(XdpAction::Drop); // Drop the paquet with the destination of port
            
        }else if is_port_in_map(transport.source() )   {
            return Ok(XdpAction::Drop);
        }
    }

    Ok(XdpAction::Pass)
}