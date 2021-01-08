#[repr(C)]
#[derive(Debug, Clone)]
pub struct PortAggs {
    pub count: u32,
}

#[repr(C)]
#[derive(Debug, Clone)]
pub struct IPAggs {
    pub count: u32,
    pub usage: u32, // bits
    // pub packet_count: u32
}
