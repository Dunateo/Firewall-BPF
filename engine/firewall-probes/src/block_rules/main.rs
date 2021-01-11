#![no_std]
#![no_main]

use firewall_probes::block_rules::{IPAggs, PortAggs};
use redbpf_probes::xdp::prelude::*;


program!(0xFFFFFFFE, "GPL");

macro_rules! insert_to_map {
    ($map:expr, $key:expr, $agg_to_insert:expr) => {
        match $map.get_mut($key) {
            Some(c) => c,
            None => {
                $map.set($key, $agg_to_insert);
                $map.get_mut($key).unwrap()
            }
        }
    }
}

#[map("ip_map")]
static mut ip_map: HashMap<u32, IPAggs> = HashMap::with_max_entries(10240);

#[map("port_map")]
static mut port_map: HashMap<u16, PortAggs> = HashMap::with_max_entries(10240);

#[xdp]
pub fn block_port_80(ctx: XdpContext) -> XdpResult {
    let ip = unsafe { *ctx.ip()? };
    let transport = ctx.transport()?;
    let data = ctx.data()?;

    let port_agg = PortAggs {
        count: 0u32
    };

    let ip_agg = IPAggs {
        count: 0u32,
        usage: 0u32, // bits
    };

    //block port 
    if let Ok(transport) = ctx.transport() {
        if transport.dest() == 80 {
            return Ok(XdpAction::Drop);

        }else if transport.source() == 80  {
            unsafe{
                let mut port_agg_sport = insert_to_map!(port_map, &transport.source(), &port_agg);
                let mut port_agg_dport = insert_to_map!(port_map, &transport.dest(), &port_agg);
                let mut ip_agg_sip = insert_to_map!(ip_map, &ip.saddr, &ip_agg);
                let mut ip_agg_dip = insert_to_map!(ip_map, &ip.daddr, &ip_agg);

                ip_agg_dip.count += 1;
                ip_agg_sip.count += 1;
                ip_agg_sip.usage += data.len() as u32 + data.offset() as u32;
                ip_agg_dip.usage += data.len() as u32 + data.offset() as u32;

                port_agg_sport.count += 1;
                port_agg_dport.count += 1;
            };
            return Ok(XdpAction::Drop);
        }
    }

    //pass to the new stack
    Ok(XdpAction::Pass)

}

