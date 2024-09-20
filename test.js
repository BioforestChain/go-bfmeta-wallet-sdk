//@ts-check
const { BCFWalletFactory } = require("@bfmeta/wallet-bcf");
(async () => {
  const bfcwallet = BCFWalletFactory({
    enable: true,
    host: [{ ip: "34.84.178.63", port: 19503 }],
    browserPath: "https://qapmapi.pmchainbox.com/browser",
  });

  console.log(bfcwallet);

  console.log(`getAddressBalance`);
  const r1 = await bfcwallet.getAddressBalance(
    "cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma",
    "XXVXQ",
    "PMC"
  );
  console.log(r1, "bfchain");

  // bfcwallet.sdk.api.transaction.broadcastCompleteTransaction();
  // bfcwallet.sdk.api.transaction.createTransferAsset();
  // bfcwallet.sdk.bfchainSignUtil.getAddressFromPublicKey
  // bfcwallet.sdk.bfchainSignUtil.createKeypair
})();
