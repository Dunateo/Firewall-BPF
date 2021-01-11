use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct BPFDataIteration {
    pub ips: Vec<EsReadyIpAggs>,
    pub ports: Vec<EsReadyPortAggs>,
}

#[derive(Deserialize, Serialize)]
pub struct EsReadyIpAggs {
    pub ip: String,
    pub count: u32,
    pub usage: u32,
}

#[derive(Deserialize, Serialize)]
pub struct EsReadyPortAggs {
    pub port: u16,
    pub count: u32,
}

fn main(){
}