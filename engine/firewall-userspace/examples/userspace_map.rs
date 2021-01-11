use serde::{Deserialize, Serialize};
use tokio;
use tokio::time::delay_for;
use redbpf::load::Loaded;
use redbpf::HashMap as BPFHashMap;
use std::time::Duration;


// aggs => aggregations
#[repr(C)]
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct PortAggs {
    pub count: u32,
}

#[repr(C)]
#[derive(Debug, Clone, Deserialize, Serialize)]
pub struct IPAggs {
    pub count: u32,
    pub usage: u32, // bits
                    // pub packet_count: u32
}

#[warn(dead_code)]
fn main(){}

