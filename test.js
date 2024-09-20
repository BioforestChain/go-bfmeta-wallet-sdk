//@ts-check
const { BCFWalletFactory } = require("@bfmeta/wallet-bcf");
const { BFMetaSignUtil } = require("@bfmeta/sign-util");
const { Buffer } = require("node:buffer");

const crypto = require("node:crypto");

class CryptoHelper {
  async sha256(msg) {
    if (msg instanceof Uint8Array) {
      return crypto.createHash("sha256").update(msg).digest();
    }
    return crypto
      .createHash("sha256")
      .update(new Uint8Array(Buffer.from(msg)))
      .digest();
  }

  async md5(data) {
    const hash = crypto.createHash("md5");
    if (data) {
      return hash.update(data).digest();
    }
    return hash;
  }

  async ripemd160(data) {
    const hash = crypto.createHash("ripemd160");
    if (data) {
      return hash.update(data).digest();
    }
    return hash;
  }
}

(async () => {
  return;
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
})();
(async () => {
  const signUtil = new BFMetaSignUtil("`+prefix+`", Buffer, new CryptoHelper());

  // bfcwallet.sdk.api.transaction.broadcastCompleteTransaction();
  // bfcwallet.sdk.api.transaction.createTransferAsset();
  // bfcwallet.sdk.bfchainSignUtil.getAddressFromPublicKey
  signUtil.createKeypair;
  signUtil.createKeypairBySecretKeyString;
  signUtil.getSecondPublicKeyStringFromSecretAndSecondSecret;
  signUtil.detachedVeriy;

  const msg = Buffer.from("utf8-1234");
  const secretKey = Buffer.from(
    "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc",
    "hex"
  );
  const publicKey = Buffer.from(
    "caf0f4c00cf9240771975e42b6672c88a832f98f01825dda6e001e2aab0bc0cc",
    "hex"
  );
  const signature = await signUtil.detachedSign(msg, secretKey);
  console.log("detachedSign", signature.toString("hex"));
  const verified = await signUtil.detachedVeriy(msg, signature, publicKey);
  console.log("verified", verified);

  // bfcwallet.sdk.bfchainSignUtil.createKeypairBySecretKeyString
})();
