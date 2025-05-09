import { crypto } from "jsr:@std/crypto";
import { encodeHex } from "jsr:@std/encoding/hex";

import $ from "@david/dax";

const sum1 = await $`cat ./go.mod | sha256sum | awk '{print $1}'`.text();

const pipe = new TransformStream<Uint8Array, Uint8Array>();
const sumP = crypto.subtle.digest("SHA-256", pipe.readable);
await $`cat ./go.mod > ${pipe.writable}`;
const sum2 = encodeHex(await sumP);

console.log(sum1);
console.log(sum2);
console.log("same? =", sum1 === sum2);

