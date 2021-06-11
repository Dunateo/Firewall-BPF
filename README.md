# NetCop
created with extendedBPF and Rust.
Choice of champions !
This project is a research Project doing a Full eBPF Personnal Firewall is not the obvious solutions but here is our research result.

The program look like this but there is much more than a simple GUI:
![GUI](https://image.noelshack.com/fichiers/2021/23/5/1623417065-capture-d-ecran-2021-06-11-a-15-10-00.png)

Our real Arhitecture it's what's working actually
![Architecture Actual](https://image.noelshack.com/fichiers/2021/23/5/1623417130-architecture.png)

The fact is when the hook tc receives the paquet the informations of the Kprobes goes to the userspace and it's already too late to compare them with user rules.
In order to filter we need to put this kind or architecture in place but RedBpf does not give us all the functionnality it will take more time:
![Architecture Progress](https://image.noelshack.com/fichiers/2021/23/5/1623417166-future-architecture.png)

Sequential architecture with the action who needs to be done in order to launch the programs

![Launching details](https://image.noelshack.com/fichiers/2021/23/5/1623417174-capture-d-ecran-2021-06-11-a-15-10-53.png)

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

`cargo build --examples and cargo make bpf`

##Running:

go in target/examples/ folder

and just run in sudo the executable

`sudo ./block_http`

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

## Thank YOU
	- RedBPF community
	- Aquarhead 
	- Junyeong Jeong
	- Daniel Borkmann
	- Kokou-milas Fokle
	- Frédéric PAILLART