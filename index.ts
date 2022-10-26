import Web3 from 'web3';
import dotenv from 'dotenv';
import { BN as BNValue } from 'bn.js';

type BN = typeof BNValue;

dotenv.config();

const SELF = process.env.WEB3_SELF_ADDRESS;
const PEER = process.env.WEB3_PEER_ADDRESS;
const KEY = process.env.WEB3_SELF_KEY;
const THEFT = process.env.WEB3_PEER_KEY;

if (!SELF || !PEER || !KEY || !THEFT) {
	console.error('Error: Environment variables not properly configured.');
	process.exit(1);
}

const web3 = new Web3(Web3.givenProvider ?? 'ws://localhost:8545');

(async () => {
	try {

		let count = 10;
		while (count--) {
			await makeTransaction(SELF, PEER, KEY, THEFT);
		}

	} catch (e: any) {
		console.error(e?.message ?? 'Error: Failed to get balance.');
		process.exit(1);
	}

	process.exit();
})();

async function makeTransaction(SELF: string, PEER: string, KEY: string, THEFT: string) {
	console.log('--- Before transaction ---');

	const selfBal = await getBal(SELF);
	console.log('Self:', selfBal, 'ETH');

	const peerBal = await getBal(PEER);
	console.log('Peer:', peerBal, 'ETH');

	console.log('--- Transaction ---');

	const spend = selfBal >= peerBal;
	const amount = Math.random() * (spend ? selfBal : peerBal);
	const to = spend ? PEER : SELF;
	// const from = spend ? SELF : PEER;
	const key = spend ? KEY : THEFT;
	const gasPrice = +await web3.eth.getGasPrice();

	console.log('Gas price:', gasPrice);

	const transaction = await web3.eth.accounts.signTransaction({
		to,
		value: toWei(amount),
		gas: 21000
	}, key);

	if (transaction.rawTransaction == null) {
		throw new Error('Failed to create transaction.');
	}

	const receipt = await web3.eth.sendSignedTransaction(transaction.rawTransaction!);
	console.log('Amount:', amount * (spend ? -1 : 1), 'ETH');
	console.log('Transaction:', receipt.transactionHash);
	console.log('Block:', receipt.blockNumber);

	console.log('--- After transaction ---');

	console.log('Self:', await getBal(SELF), 'ETH');
	console.log('Peer:', await getBal(PEER), 'ETH');
}

async function getBal(add: string): Promise<number> {
	return +web3.utils.fromWei(await web3.eth.getBalance(add), 'ether');
}

function toWei(val: number | BigInt | string | BN): string {
	let paramVal: string | number;

	if (typeof val === 'bigint')
		paramVal = val.toString();
	else
		paramVal = val as (string | number);

	return web3.utils.toWei(paramVal.toString(), 'ether').toString();
}