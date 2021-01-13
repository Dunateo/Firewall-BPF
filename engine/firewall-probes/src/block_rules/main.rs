#![no_std]
#![no_main]

use redbpf_probes::tc::*;

program!(0xFFFFFFFE, "GPL");

#[tc_action]
pub fn block_port_80(ctx: TcActionResult) -> TcActionResult {
    if let Ok(transport) = ctx.transport() {
        if transport.dest() == 80 {
            return Ok(TcAction::Shot);
            
        }else if transport.source() == 80  {
            return Ok(TcAction::Shot);
        }
    }

    Ok(TcAction::Ok)
}