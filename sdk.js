async function returnToGo(req_id, handler) {
    console.log("node env ready", process.versions.node)
    try {
        const result = await handler()
        return `Result ${req_id} ${JSON.stringify(result)}`
    } catch (e) {
        return `Result ${req_id} ${String(e)}`
    }
}

Object.assign(globalThis, {returnToGo});
console.log("node env ready", process.versions.node)
const repl = require('node:repl');
repl.start({
    prompt: '',
    writer: (output) => {
        console.log(output)
    }
});

