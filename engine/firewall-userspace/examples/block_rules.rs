use std::process;
use redbpf::load::Loader;

fn main() -> Result<()>{
    if unsafe { libc::getuid() } != 0 {
        eprintln!("You must be root to use eBPF!");
        process::exit(1);
    }


        //let mut module = Module::parse(&std::fs::read("vfsreadlat.elf").unwrap()).unwrap();
        let mut module = Loader::load(probe_code()).expect("error loading BPF program");
        //let mut module = Loader::load_file("vfsreadlat.elf").expect("error loading probe");
            for prog in module.module.programs.iter_mut() {
                prog.load(module.module.version, module.module.license.clone())
                    .expect("failed to load program");
            }


    loop {}
}

fn probe_code() -> &'static [u8] {
    include_bytes!(concat!(
        env!("OUT_DIR"),
        "/target/bpf/programs/block_rules/block_rules.elf"
    ))
}
