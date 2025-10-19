// export async function hashSha512(data: Uint8Array): Promise<string> {
//   const hashBuffer = await window.crypto.subtle.digest("SHA-512", data);
//   const hashArray = Array.from(new Uint8Array(hashBuffer));
//   const hashHex = hashArray
//     .map((b) => b.toString(16).padStart(2, "0"))
//     .join(""); // convert bytes to hex string
//   return hashHex;
// }

export async function hashSha512(data: Uint8Array): Promise<string> {
  const hashBuffer = await crypto.subtle.digest("SHA-512", data.buffer as ArrayBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}
