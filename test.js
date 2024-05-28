(async () => {
    const bfcwallet = new require('@bfmeta/wallet-bcf').BCFWalletFactory({
        enable: true, host: [{ip: "34.84.178.63", port: 19503}], browserPath: "https://qapmapi.pmchainbox.com/browser",
    });

    console.log(bfcwallet)

    console.log(`getAddressBalance`);
    const r1 = await bfcwallet.getAddressBalance("cEAXDkaEJgWKMM61KYz2dYU1RfuxbB8Ma", "XXVXQ", "PMC");
    console.log(r1, "bfchain");
})();
