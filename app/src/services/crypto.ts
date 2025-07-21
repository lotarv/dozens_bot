export async function decryptGoAES( encrypted: string, password: string): Promise<string> {
    const [b64Salt, b64Nonce, b64Ciphertext] = encrypted.split(':');
    if (!b64Salt || !b64Nonce || !b64Ciphertext) {
        throw new Error('Invalid encrypted format');
    }

    const salt = Uint8Array.from(atob(b64Salt), c => c.charCodeAt(0));
    const nonce = Uint8Array.from(atob(b64Nonce), c => c.charCodeAt(0));
    const ciphertext = Uint8Array.from(atob(b64Ciphertext), c => c.charCodeAt(0));

    const keyMaterial = await crypto.subtle.importKey(
        'raw',
        new TextEncoder().encode(password),
        'PBKDF2',
        false,
        ['deriveKey']
    );

    const key = await crypto.subtle.deriveKey(
        {
        name: 'PBKDF2',
        salt,
        iterations: 100_000,
        hash: 'SHA-256',
        },
        keyMaterial,
        { name: 'AES-GCM', length: 256 },
        false,
        ['decrypt']
    );

    const decrypted = await crypto.subtle.decrypt(
        {
        name: 'AES-GCM',
        iv: nonce,
        },
        key,
        ciphertext
    );

    return new TextDecoder().decode(decrypted);
}

export function isEncrypted(text: string): boolean {
  return typeof text === "string" && text.split(":").length === 3
}