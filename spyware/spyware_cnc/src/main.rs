mod command;
mod server;
use server::run_server;

#[macro_use]
extern crate failure;

pub const SERVER_BIND_PORT: u16 = 9393;

fn main() {
    // Blocking call
    run_server(SERVER_BIND_PORT);
    panic!("Server exited unexpectedly");
}
