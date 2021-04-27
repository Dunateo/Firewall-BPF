use anyhow::Result;
use redbpf::load::{Loaded, Loader};

fn main() -> Result<()> {
  let prog = include_bytes!("/target/bpf/programs/limit/limit.elf");
  //let mut module = Module::parse(prog).expect("error parsing BPF code");
  let mut module = Loader::load(prog).expect("error loading BPF program");

  for program in module.module.programs.iter_mut() {
    program
      .load(module.module.version, module.module.license.clone())
      .expect("failed to load program");
  }

  loop {}
}
