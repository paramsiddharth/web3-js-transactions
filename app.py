from __future__ import annotations
from dotenv import load_dotenv

load_dotenv()

from os import environ as env
from sys import exit, stderr
from random import random
from web3 import Web3, EthereumTesterProvider
from web3.types import Wei

SELF  = env.get('WEB3_SELF_ADDRESS')
PEER  = env.get('WEB3_PEER_ADDRESS')
KEY   = env.get('WEB3_SELF_KEY')
THEFT = env.get('WEB3_PEER_KEY')

if None in (SELF, PEER, KEY, THEFT):
	print('Error: Environment variables not properly configured.', file=stderr)
	exit(1)

url = env.get('WEB3_PROVIDER', 'ws://localhost:8545')
w3 = Web3(Web3.WebsocketProvider(url) if url else EthereumTesterProvider())

def main():
	try:

		count = 10
		for i in range(count):
			i += 1
			print('\n--- Iteration', i, '---\n')
			make_transaction(SELF, PEER, KEY, THEFT)

	except Exception as e:
		print(', '.join(e.args), file=stderr)
		exit(1)

	exit()

def make_transaction(
	SELF: str,
	PEER: str,
	KEY: str,
	THEFT: str
):
	print('--- Before transaction ---')

	self_bal = get_bal(SELF)
	print('Self:', self_bal, 'ETH')

	peer_bal = get_bal(PEER)
	print('Peer:', peer_bal, 'ETH')

	print('--- Transaction ---')

	spend = self_bal >= peer_bal
	amount = random() * float(self_bal if spend else peer_bal)
	to = PEER if spend else SELF
	# from = SELF if spend else PEER
	key = KEY if spend else THEFT
	gas_price = w3.fromWei(w3.eth.gasPrice, 'ether')

	print('Gas price:', gas_price, 'ETH')

	transaction = w3.eth.account.sign_transaction({
		'to': to,
		'value': to_wei(amount),
		'gas': 21_000,
		'gasPrice': w3.eth.gasPrice,
		'nonce': 0
	}, key)

	hex = w3.eth.send_raw_transaction(transaction.rawTransaction)
	receipt = w3.eth.get_transaction_receipt(hex)

	print('Amount:', amount * (-1 if spend else 1), 'ETH')
	print('Transaction:', hex.hex())
	print('Block:', receipt.blockNumber)

	print('--- After transaction ---')

	print('Self:', get_bal(SELF), 'ETH')
	print('Peer:', get_bal(PEER), 'ETH')

def get_bal(add: str) -> float:
	return w3.fromWei(w3.eth.get_balance(add), 'ether')

def to_wei(val: int | float | str) -> Wei:
	return w3.toWei(val, 'ether')

if __name__ == '__main__':
	main()