#![allow(non_snake_case)]

use std::{env, fs, process};
use rand::{self, Rng, distributions::uniform::SampleRange};

use dotenv;
use web3::{transports, Web3, types::{U256, H160}, Transport};

#[tokio::main]
async fn main() -> web3::Result<()> {
	let env_dir = env::current_dir()
		.and_then(|d| Ok(d.join("../.env")))?;
	let env_dir = fs::canonicalize(&env_dir).unwrap();
	dotenv::from_path(env_dir.as_path()).ok();

	let SELF = env::var("WEB3_SELF_ADDRESS").unwrap_or_else(|_| String::from(""));
	let PEER = env::var("WEB3_PEER_ADDRESS").unwrap_or_else(|_| String::from(""));
	let KEY = env::var("WEB3_SELF_KEY").unwrap_or_else(|_| String::from(""));
	let THEFT = env::var("WEB3_PEER_KEY").unwrap_or_else(|_| String::from(""));

	if SELF.len() < 1
		|| PEER.len() < 1
		|| KEY.len() < 1
		|| THEFT.len() < 1 {
		eprintln!("Error: Environment variables not properly configured.");
		process::exit(1);
	}

	let PROVIDER = env::var("WEB3_PROVIDER").unwrap_or_else(|_| String::from(""));

	let transport = transports::WebSocket::new(
		if !PROVIDER.is_empty() { PROVIDER.as_str() } else { "ws://localhost:8545" }
	).await?;
	let web3 = Web3::new(transport);

	let count = 10;
	for i in 1..=count {
		println!("\n--- Iteration {i} ---\n");
		make_transaction(&web3, &SELF, &PEER, &KEY, &THEFT).await;
	}

	Ok(())
}

async fn make_transaction<T: Transport>(
	web3: &Web3<T>,
	SELF: &str,
	PEER: &str,
	KEY: &str,
	THEFT: &str
) -> Option<()> {
	let mut rng = rand::thread_rng();

	println!("--- Before transaction ---");

	let self_bal = get_bal(web3, SELF).await.unwrap();
	println!("Self: {self_bal} ETH");

	let peer_bal = get_bal(web3, PEER).await.unwrap();
	println!("Peer: {peer_bal} ETH");

	println!("--- Transaction ---");

	let spend = self_bal >= peer_bal;
	let amount = rng.gen_range(0..1) as f64 * if spend { self_bal } else { peer_bal };
	let to = if spend { PEER } else { SELF };
	// let from = if spend { SELF } else { PEER };
	let key  = if spend { KEY } else { THEFT };
	let gas_price = web3.eth()
		.gas_price()
		.await
		.unwrap()
		.low_u128()
		as f64
		/ 1e18;

	println!("Gas price: {gas_price} ETH");

	// let transaction = web3.eth()

	println!("--- After transaction ---");

	println!("Self: {} ETH", get_bal(web3, SELF).await.unwrap());
	println!("Peer: {} ETH", get_bal(web3, PEER).await.unwrap());

	Some(())
}

async fn get_bal<T: Transport>(web3: &Web3<T>, add: &str) -> Result<f64, web3::Error> {

	return Ok(web3.eth().balance(
		add.parse::<H160>().unwrap(), None)
		.await?
		.low_u128()
		as f64
		/ 1e18
	);
}

fn to_wei(val: f64) -> U256 {
	U256::from((val * 1e18).ceil() as u128)
}