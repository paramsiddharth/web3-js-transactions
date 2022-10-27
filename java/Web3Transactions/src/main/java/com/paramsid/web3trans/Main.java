package com.paramsid.web3trans;

import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;

import org.apache.commons.lang3.StringUtils;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.RawTransaction;
import org.web3j.crypto.TransactionEncoder;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.http.HttpService;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;

import io.github.cdimascio.dotenv.Dotenv;

public class Main
{
	static Web3j web3 = null;

	public static void main( String[] args )
	{
		Dotenv env = Dotenv.configure().directory("../..").load();

		final String PROVIDER  = env.get("WEB3_PROVIDER");

		final String SELF  = env.get("WEB3_SELF_ADDRESS");
		final String PEER  = env.get("WEB3_PEER_ADDRESS");
		final String KEY   = env.get("WEB3_SELF_KEY");
		final String THEFT = env.get("WEB3_PEER_KEY");

		if (
			   StringUtils.isEmpty(SELF)
			|| StringUtils.isEmpty(PEER)
			|| StringUtils.isEmpty(KEY)
			|| StringUtils.isEmpty(THEFT)
		) {
			System.err.println("Error: Environment variables not properly configured.");
			System.exit(1);
		}

		web3 = Web3j.build(!StringUtils.isEmpty(PROVIDER) ? new HttpService(PROVIDER) : new HttpService());

		try {

			var count = 10;
			for (int i = 1; i <= count; i++) {
				System.out.println("\n--- Iteration " + i + " ---\n");
				makeTransaction(SELF, PEER, KEY, THEFT);
			}

		} catch (Exception e) {
			System.err.println(e.getMessage());
			System.exit(1);
		}

		System.exit(0);
	}

	static void makeTransaction(
		final String SELF,
		final String PEER,
		final String KEY,
		final String THEFT
	) throws IOException {
		assert web3 != null;

		System.out.println("--- Before transaction ---");
		var selfBal = getBal(SELF);
		System.out.println("Self: " + selfBal + " ETH");

		var peerBal = getBal(PEER);
		System.out.println("Peer: " + peerBal + " ETH");

		System.out.println("--- Transaction ---");

		var spend = selfBal >= peerBal;
		var amount = Math.random() * (spend ? selfBal : peerBal);
		var to = spend ? PEER : SELF;
		// var from = spend ? SELF : PEER;
		var key = spend ? KEY : THEFT;
		var gasPrice = Convert.fromWei(web3.ethGasPrice().send().getGasPrice().toString(), Convert.Unit.ETHER);

		System.out.println("Gas price: " + gasPrice + " ETH");

		var transaction = RawTransaction.createEtherTransaction(BigInteger.valueOf(0), web3.ethGasPrice().send().getGasPrice(), BigInteger.valueOf(21000), to, toWei(amount));
		var creds = Credentials.create(key);
		var signed = TransactionEncoder.signMessage(transaction, creds);
		var receipt = web3.ethSendRawTransaction(Numeric.toHexString(signed)).send();
		String hash = receipt.getTransactionHash();

		System.out.println("Amount: " + amount * (spend ? -1 : 1) + " ETH");
		System.out.println("Transaction: " + hash);
		System.out.println("Block: " + web3.ethGetTransactionByHash(hash).send().getTransaction().get().getBlockNumber());

		System.out.println("--- After transaction ---");

		System.out.println("Self: " + getBal(SELF) + " ETH");
		System.out.println("Peer: " + getBal(PEER) + " ETH");
	}

	static double getBal(String add) throws IOException {
		assert web3 != null;
		var recv = web3.ethGetBalance(add, DefaultBlockParameterName.LATEST).send().getBalance();
		var balance = Convert.fromWei(recv.toString(), Convert.Unit.ETHER);
		return balance.doubleValue();
	}

	static BigInteger toWei(double val) {
		return Convert.toWei(BigDecimal.valueOf(val), Convert.Unit.ETHER).toBigInteger();
	}
}