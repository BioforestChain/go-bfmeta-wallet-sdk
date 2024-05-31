#!/usr/bin/env node
if (process.argv.includes("--version")) {
    const pkg = require("./package.json");
    console.log(pkg.name, pkg.version);
    process.exit(0);
}
const crypto = require("node:crypto")

class CryptoHelper {
    async sha256(msg) {
        if (msg instanceof Uint8Array) {
            return crypto
                .createHash("sha256")
                .update(msg)
                .digest()
        }
        return crypto
            .createHash("sha256")
            .update(new Uint8Array(Buffer.from(msg)))
            .digest()
    }

    async md5(data) {
        const hash = crypto.createHash("md5")
        if (data) {
            return hash.update(data).digest()
        }
        return hash
    }

    async ripemd160(data) {
        const hash = crypto.createHash("ripemd160")
        if (data) {
            return hash.update(data).digest()
        }
        return hash
    }
}

globalThis.cryptoHelper = new CryptoHelper();
globalThis.walletBcf = require("@bfmeta/wallet-bcf");
globalThis.__signUtil = require("@bfmeta/sign-util");

async function returnToGo(req_id, handler) {
    console.log("node env ready", process.versions.node);
    try {
        const result = await handler();
        return `Result ${req_id} ${JSON.stringify(result)}`;
    } catch (e) {
        return `Result ${req_id} ${String(e)}`;
    }
}

Object.assign(globalThis, {returnToGo});
console.log("node env ready", process.versions.node);
const repl = require("node:repl");
repl.start({
    prompt: "",
    writer: (output) => {
        console.log(output);
    },
});
