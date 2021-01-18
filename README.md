# NetCop
created with extendedBPF and Rust.
Choice of champions !

![Architecture Front](https://image.noelshack.com/fichiers/2021/03/1/1610990015-architecture-bpf-min.png)
![Architecture Back](https://image.noelshack.com/fichiers/2021/03/1/1610990015-architecture-bpf-min.png)

## Installation

# llvm-11

[https://apt.llvm.org/](https://apt.llvm.org/)

`bash -c "$(wget -O - [https://apt.llvm.org/llvm.sh](https://apt.llvm.org/llvm.sh))`

lvm-config --version **should work**

# Install dependencies

Preparation for the cargo installation.

```bash
sudo apt-get -y install build-essential zlib1g-dev \
		llvm-11-dev libclang-11-dev linux-headers-$(uname -r)
```

# Install Rust

`curl --proto '=https' --tlsv1.2 -sSf [https://sh.rustup.rs](https://sh.rustup.rs/) | sh`

Don't forget to copy the environment PATH

# Examples

```bash
cargo install cargo-bpf
```

You might have some errors.

#Compilation:

`cargo build --examples`

##Running:

go in target/examples/ folder

and just run in sudo the executable

`sudo ./block_http`

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
